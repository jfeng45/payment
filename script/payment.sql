CREATE DATABASE `service_config` ;

/*Table structure for table `userinfo` */

DROP TABLE IF EXISTS `payment`;

CREATE TABLE `payment` (
  `id` int(10) NOT NULL auto_increment,
  `source_account` int(10) default 0,
  `target_account` int(10) default 0,
  `amount` decimal(12,2) default 0.00,
  `payment_method` varchar(16) default NULL,
  `status` varchar(16) NULL,
  `order_number` varchar(16) default NULL,
  `created_time` DATETIME NOT NULL,
  `completion_time` DATETIME DEFAULT NULL,
  PRIMARY KEY  (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
