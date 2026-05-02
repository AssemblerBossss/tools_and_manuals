// Copyright © 2022 Aleksey Barabanov
// Licensed under the Apache License, Version 2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/streadway/amqp"
)

var (
	rabbitFromUser       = GetEnv("FROM_USERNAME", "rmuser")
	rabbitFromPassword   = url.QueryEscape(GetEnv("FROM_PASSWORD", "rmpassword"))
	rabbitFromHost       = GetEnv("FROM_HOSTNAME", "rabbitmq")
	rabbitFromPort       = GetEnv("FROM_PORT", "5672")
	rabbitFromQueue      = GetEnv("FROM_QUEUE", "")
	rabbitFromExchange   = GetEnv("FROM_EXCHANGE", "")
	rabbitFromRoutingKey = GetEnv("FROM_ROUTINGKEY", "")
	rabbitFromPassive    = GetEnv("FROM_PASSIVE", "true")
	rabbitPrefetch       = GetEnv("PREFETCH", "5")
	fail                 = GetEnv("FAIL", "false")
	rejectall            = GetEnv("REJECTALL", "false")
	sleepms              = convertTextInt(GetEnv("SLEEP", "0"))
	retryQueue           = GetEnvBool("RETRY_QUEUE", false)
	retryQueueTtl        = GetEnvInt("RETRY_QUEUE_TTL", 15000)

	uriFrom      = flag.String("uriFrom", "amqp://"+rabbitFromUser+":"+rabbitFromPassword+"@"+rabbitFromHost+":"+rabbitFromPort+"/", "AMQP From URI")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	consumerTag  = flag.String("consumer-tag", os.Getenv("HOSTNAME"), "AMQP consumer tag (should not be blank)")
)

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func init() {
	flag.Parse()
}

func main() {
	// Запуск консьюмера
	c, err := newConsumer()
	CheckErr(err)

	// Консьюмер запускается фоном - поэтому делаем ожидание останова приложения для останова всего демона
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.
	// Закрываем канал
	c.channel.Cancel(c.tag, true)
	// Закрываем соединение
	c.conn.Close()
}

// Consumer запуск консьюмера
func newConsumer() (*consumer, error) {
	c := &consumer{
		conn:    nil,
		channel: nil,
		tag:     *consumerTag,
		done:    make(chan error),
	}

	var err error
	// Создаём соединение
	// Просто по умолчанию (10 секунд heartbeat, таймаут 30 сек)
	// c.conn, err = amqp.Dial(*uriFrom)

	// Явно обозначаем значение heartbeat
	config := amqp.Config{
		Heartbeat: time.Second * 12,
	}

	log.Printf("FROM: dialing")
	// Соединение
	c.conn, err = amqp.DialConfig(*uriFrom, config)

	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}
	log.Printf("FROM: got Connection, getting Channel")
	// Создаём канал
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("FROM: Channel: %s", err)
	}

	if rabbitFromPassive != "true" { // Если не пэссив
		// Если необходима очередь повторных попыток
		queueArg := amqp.Table{}
		if retryQueue {
			queueArg = amqp.Table{"x-dead-letter-exchange": rabbitFromQueue + ".fail"}
			retryQueueArg := amqp.Table{
				"x-dead-letter-exchange": rabbitFromQueue + ".retry",
				"x-message-ttl":          retryQueueTtl}
			log.Printf("FROM: declaring Retry Queue (%s)", rabbitFromQueue+".fail")
			_, err := c.channel.QueueDeclare(
				rabbitFromQueue+".retry.dlx", // name of the queue
				true,                         // durable
				false,                        // delete when usused
				false,                        // exclusive
				false,                        // noWait
				retryQueueArg,                // arguments
			)
			CheckErrLog(err, "Retry Queue declare fail")
			err = c.channel.ExchangeDeclare(
				rabbitFromQueue+".fail", // name of the exchange
				"direct",                // type
				true,                    // durable
				false,                   // delete when complete
				false,                   // internal
				false,                   // noWait
				nil,                     // arguments
			)
			CheckErrLog(err, "Retry Queue exchange declare fail")
			err = c.channel.ExchangeDeclare(
				rabbitFromQueue+".retry", // name of the exchange
				"direct",                 // type
				true,                     // durable
				false,                    // delete when complete
				false,                    // internal
				false,                    // noWait
				nil,                      // arguments
			)
			CheckErrLog(err, "Retry Queue exchange declare fail")
			err = c.channel.QueueBind(
				rabbitFromQueue+".retry.dlx", // name of the queue
				rabbitFromRoutingKey,         // bindingKey
				rabbitFromQueue+".fail",      // sourceExchange
				false,                        // noWait
				nil,                          // arguments
			)
			CheckErrLog(err, "Retry Queue binding declare fail")
			err = c.channel.QueueBind(
				rabbitFromQueue,          // name of the queue
				rabbitFromRoutingKey,     // bindingKey
				rabbitFromQueue+".retry", // sourceExchange
				false,                    // noWait
				nil,                      // arguments
			)
			CheckErrLog(err, "Retry Queue binding declare fail")
		}
		// Декларируем очередь
		log.Printf("FROM: declaring Queue (%s)", rabbitFromQueue)
		_, err := c.channel.QueueDeclare(
			rabbitFromQueue, // name of the queue
			true,            // durable
			false,           // delete when usused
			false,           // exclusive
			false,           // noWait
			queueArg,        // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("queue Declare: %s", err)
		}
		if rabbitFromExchange != "" && rabbitFromRoutingKey != "" {
			log.Printf("FROM: got Channel, declaring Exchange (%s)", rabbitFromExchange)
			// Декларируем эксчейндж
			if err = c.channel.ExchangeDeclare(
				rabbitFromExchange, // name of the exchange
				*exchangeType,      // type
				true,               // durable
				false,              // delete when complete
				false,              // internal
				false,              // noWait
				nil,                // arguments
			); err != nil {
				return nil, fmt.Errorf("exchange Declare: %s", err)
			}
			log.Printf("FROM: binding Exchange (%s) to Queue (%s) with RoutingKey (%s)", rabbitFromExchange, rabbitFromQueue, rabbitFromRoutingKey)
			// Декларируем биндинг
			if err = c.channel.QueueBind(
				rabbitFromQueue,      // name of the queue
				rabbitFromRoutingKey, // bindingKey
				rabbitFromExchange,   // sourceExchange
				false,                // noWait
				nil,                  // arguments
			); err != nil {
				return nil, fmt.Errorf("queue Bind: %s", err)
			}
		}
	} else { // Passive mode
		log.Printf("FROM: passive declaring Queue (%s)", rabbitFromQueue)
		// Декларируем очередь пассивно (если нет то упадет с ошибкой)
		_, err := c.channel.QueueDeclarePassive(
			rabbitFromQueue, // name of the queue
			true,            // durable
			false,           // delete when usused
			false,           // exclusive
			false,           // noWait
			nil,             // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("queue Declare: %s", err)
		}
		if rabbitFromExchange != "" && rabbitFromRoutingKey != "" {
			log.Printf("FROM: got Channel, passive declaring Exchange (%s)", rabbitFromExchange)
			// Декларируем эксчейндж пассивно (если нет то упадет с ошибкой)
			if err = c.channel.ExchangeDeclarePassive(
				rabbitFromExchange, // name of the exchange
				*exchangeType,      // type
				true,               // durable
				false,              // delete when complete
				false,              // internal
				false,              // noWait
				nil,                // arguments
			); err != nil {
				return nil, fmt.Errorf("exchange Declare: %s", err)
			}
			log.Printf("FROM: binding Exchange (%s) to Queue (%s) with RoutingKey (%s)", rabbitFromExchange, rabbitFromQueue, rabbitFromRoutingKey)
			// Декларируем биндинг
			if err = c.channel.QueueBind(
				rabbitFromQueue,      // name of the queue
				rabbitFromRoutingKey, // bindingKey
				rabbitFromExchange,   // sourceExchange
				false,                // noWait
				nil,                  // arguments
			); err != nil {
				return nil, fmt.Errorf("queue Bind: %s", err)
			}
		}
	}

	// Если не обьявлена очередь - не пытаться считать ее длину
	if rabbitFromQueue != "" {
		// Через канал можно получать всякую информацию
		state, err := c.channel.QueueInspect(rabbitFromQueue)
		CheckErrLog(err, "Get state from queue")
		// Например о длине очереди
		log.Printf("FROM: %d messages in queue %v\n", state.Messages, rabbitFromQueue)
	}

	// Преобразуем префеч из строки в int
	prcount, err := strconv.Atoi(rabbitPrefetch)
	if err != nil {
		return nil, fmt.Errorf("unable to convert preFetchCount: %s", err)
	}
	// Устанавливаем префетч каунт канала
	if err = c.channel.Qos(
		prcount,
		0,
		false,
	); err != nil {
		return nil, fmt.Errorf("error qos: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag '%s'), sleeping: %v ms", c.tag, sleepms)
	deliveries, err := c.channel.Consume(
		rabbitFromQueue, // name
		c.tag,           // consumerTag,
		false,           // autoAck
		false,           // exclusive
		false,           // noLocal
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %s", err)
	}

	// Стартуем рутину с хендлером сообщений
	go handle(deliveries, c.done)

	return c, nil
}

// handle обработчик сообщений консьюмера.
func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for msg := range deliveries {
		// Выполняется в момент получения сообщения
		if sleepms != 0 {
			time.Sleep(time.Millisecond * time.Duration(sleepms)) // задержка
		}
		if rejectall == "true" {
			msg.Reject(false)
		} else if fail == "true" {
			random := rand.Intn(100)
			log.Printf("Message: %s, RoutingKey: %s, random: %v", string(msg.Body), msg.RoutingKey, random)
			if random >= 50 {
				msg.Ack(false)
			} else {
				log.Printf("FAIL - reject")
				msg.Reject(false)
			}
		} else {
			log.Printf("Message: %s, RoutingKey: %s", string(msg.Body), msg.RoutingKey)
			msg.Ack(false)
		}

	}
	// При разрушении канала выключаемся
	log.Fatalln("handle: deliveries channel closed, shutting down")
	done <- nil
}

// GetEnv Получение данных из переменных окружения, также передаем дефолт если переменной не найдется
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvBool Получение данных из переменных окружения, также передаем дефолт если переменной не найдется - булево
func GetEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if value == "true" {
			return true
		} else {
			return false
		}
	}
	return fallback
}

// GetEnvInt Получение данных из переменных окружения, также передаем дефолт если переменной не найдется - число
func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		result, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("unable to convert %v: %s", key, err)
		}
		return result
	}
	return fallback
}

// CheckErr Проверка ошибки
func CheckErr(err error) { // Check error
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// CheckErrLog Проверка ошибки со строковым префиксом (комментом)
func CheckErrLog(err error, comment string) { // Check error
	if err != nil {
		log.Fatalf("%v : %v", comment, err)
	}
}

// convertTextInt Конвертация строки в число
func convertTextInt(str string) int {
	result, err := strconv.Atoi(str)
	CheckErrLog(err, "error converting")
	return result
}
