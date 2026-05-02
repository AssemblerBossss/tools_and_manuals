// Copyright © 2022 Aleksey Barabanov
// Licensed under the Apache License, Version 2.0

package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"net/url"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

var (
	rabbitToUser       = GetEnv("TO_USERNAME", "rmuser")
	rabbitToPassword   = url.QueryEscape(GetEnv("TO_PASSWORD", "rmpassword"))
	rabbitToHost       = GetEnv("TO_HOSTNAME", "rabbitmq")
	rabbitToPort       = GetEnv("TO_PORT", "5672")
	rabbitToQueue      = GetEnv("TO_QUEUE", "")
	rabbitToExchange   = GetEnv("TO_EXCHANGE", "")
	rabbitToRoutingKey = GetEnv("TO_ROUTINGKEY", "")
	rabbitToPassive    = GetEnv("TO_PASSIVE", "true")
	mchMode            = GetEnv("MCH_MODE", "false")
	jobMode            = GetEnv("JOB_MODE", "false")
	sleepms            = convertTextInt(GetEnv("SLEEP", "0"))
	report             = convertTextInt(GetEnv("REPORT", "10000"))
	genCount           = convertTextInt(GetEnv("COUNT", "100000"))

	envDebug = GetEnv("APP_DEBUG", "false")
	// AppDebug Debug log application
	AppDebug     bool
	uriTo        = flag.String("uriTo", "amqp://"+rabbitToUser+":"+rabbitToPassword+"@"+rabbitToHost+":"+rabbitToPort+"/", "AMQP To URI")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
)

func init() {
	flag.Parse()
}

func main() {
	AppDebug = false
	if envDebug == "true" {
		AppDebug = true
		log.Printf("DEBUG ON !")
	}

	log.Print("starting ...")
	amqpConnectTo, channelTo, err := connectTo()
	CheckErrLog(err, "ConnectTo")

	if mchMode == "true" {
		log.Print("starting in Multi Channel Mode")
		genChannles(amqpConnectTo)
	} else {
		log.Print("starting in Single Channel Mode...")
		gen(channelTo)
	}

	log.Printf("close connection")
	err = amqpConnectTo.Close()
	CheckErrLog(err, "close connection error")
	if jobMode != "true" {
		log.Printf("while true...")
		select {}
	}
}

// connectTo - функция создания подключения к реббиту назначения, возвращает ссылку на обьект коннекта и ссылку на обьект канала
func connectTo() (*amqp.Connection, *amqp.Channel, error) {
	// Подключаемся к реббиту
	if AppDebug {
		log.Printf("TO: dialing %s", *uriTo)
	}
	connection, err := amqp.Dial(*uriTo)
	CheckErrLog(err, "Connect to rabbit_to")
	// defer connection.Close()

	if AppDebug {
		log.Printf("TO: got Connection, getting Channel")
	}
	channel, err := connection.Channel()
	CheckErrLog(err, "Initial getChannel TO")

	// Ежели не указан эксчейндж создаем очередь одноименную с роутинг кеем
	if rabbitToPassive != "true" { // Если не пассивный режим
		if rabbitToQueue != "" {
			log.Printf("TO: declaring Queue (%s)", rabbitToQueue)
			_, err := channel.QueueDeclare(
				rabbitToQueue, // name of the queue
				true,          // durable
				false,         // delete when usused
				false,         // exclusive
				false,         // noWait
				nil,           // arguments
			)
			CheckErrLog(err, "TO declaring Queue")
		}
		if rabbitToExchange != "" { // Если есть название эксчейнджа - декларируем и его
			log.Printf("TO: declaring Exchange (%s)", rabbitToExchange)
			err = channel.ExchangeDeclare(
				rabbitToExchange, // name of the exchange
				*exchangeType,    // type
				true,             // durable
				false,            // delete when complete
				false,            // internal
				false,            // noWait
				nil,              // arguments
			)
			CheckErr(err)
		}
		if rabbitToExchange != "" && rabbitToQueue != "" {
			log.Printf("TO: binding Exchange (%s) to Queue (%s) with RoutingKey (%s)", rabbitToExchange, rabbitToQueue, rabbitToRoutingKey)
			err = channel.QueueBind(
				rabbitToQueue,      // name of the queue
				rabbitToRoutingKey, // bindingKey
				rabbitToExchange,   // sourceExchange
				false,              // noWait
				nil,                // arguments
			)
			CheckErr(err)
		}

	} else { // Пассивный режим
		if rabbitToQueue != "" {
			log.Printf("TO: passive declaring Queue (%s)", rabbitToQueue)
			_, err := channel.QueueDeclarePassive(
				rabbitToQueue, // name of the queue
				true,          // durable
				false,         // delete when usused
				false,         // exclusive
				false,         // noWait
				nil,           // arguments
			)
			CheckErr(err)
		}
		if rabbitToExchange != "" { // Если есть название эксчейнджа - декларируем и его
			log.Printf("TO: passive declaring Exchange (%s)", rabbitToExchange)
			err = channel.ExchangeDeclarePassive(
				rabbitToExchange, // name of the exchange
				*exchangeType,    // type
				true,             // durable
				false,            // delete when complete
				false,            // internal
				false,            // noWait
				nil,              // arguments
			)
			CheckErr(err)

		}
		if rabbitToExchange != "" && rabbitToQueue != "" {
			log.Printf("TO: binding Exchange (%s) to Queue (%s) with RoutingKey (%s)", rabbitToExchange, rabbitToQueue, rabbitToRoutingKey)
			err = channel.QueueBind(
				rabbitToQueue,      // name of the queue
				rabbitToRoutingKey, // bindingKey
				rabbitToExchange,   // sourceExchange
				false,              // noWait
				nil,                // arguments
			)
			CheckErr(err)
		}
	}

	// Если не обьявлена очередь - не пытаться считать ее длину
	if rabbitToQueue != "" {
		state, err := channel.QueueInspect(rabbitToQueue)
		CheckErrLog(err, "Final state check channel TO")
		fmt.Printf("%d messages in queue\n", state.Messages)
	}

	return connection, channel, nil
}

func gen(channel *amqp.Channel) {
	for count := 1; count <= genCount; count++ {
		byteMessage := []byte(strconv.Itoa(count))
		publish(channel, byteMessage)
		// log.Printf("Generated %v\n", count)
		if count%report == 0 {
			// fmt.Printf("%v receipts\n", count)
			log.Printf("Delivered %v messages\n", count)
		}
	}
}

func genChannles(connection *amqp.Connection) {
	for count := 1; count <= genCount; count++ {
		byteMessage := []byte(strconv.Itoa(count))
		channel, err := connection.Channel()
		CheckErr(err)
		publish(channel, byteMessage)
		channel.Close()
		// log.Printf("Generated %v\n", count)
		if count%10000 == 0 {
			// fmt.Printf("%v receipts\n", count)
			log.Printf("Delivered %v messages\n", count)
		}
	}
}

// publish метод отправки сообщения
func publish(channel *amqp.Channel, body []byte) error {
	// log.Printf("Generated %v\n", body)
	routing := rabbitToRoutingKey
	if rabbitToExchange == "" {
		routing = rabbitToQueue
	}

	err := channel.Publish(
		rabbitToExchange, // publish to an exchange
		routing,          // routing to 0 or more queues
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    2, // 1=non-persistent, 2=persistent
			Priority:        0, // 0-9
			// a bunch of application/implementation-specific fields
		},
	)
	time.Sleep(time.Millisecond * time.Duration(sleepms))
	CheckErrLog(err, "error publishing")
	return nil

}
