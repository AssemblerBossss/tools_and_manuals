<?php
 
/*
 * This is an example of how to call the Erlang class
 */
 
require ("erlang2.class.php");
 
$erlang = new Erlang();
 
//Debugging is disabled by default, and are the error messages. If these are needed, they can be enabled below
// $erlang->set_silent(0);
$erlang->set_debug(0);

$calls=		$_POST['calls'];
$interval=	$_POST['interval'];
$maxq=	$_POST['maxq'];
$ctime=		$_POST['ctime'];
$atime=		$_POST['atime'];
$uo=		$_POST['uo'];
 

$erlang->set_parameters($calls,$ctime,$maxq,$interval);
echo "Consumers count: ".$erlang->calculate_required_agents($atime,$uo) . "<hr><br>";

$erlang->set_debug(1);
print_r($_POST);
// $erlang->set_parameters($calls,$ctime,$maxq,$interval);
echo $erlang->calculate_required_agents($atime,$uo) . "<br>";
