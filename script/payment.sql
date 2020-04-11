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
  `created_time` date default NULL,
  `completion_time` date default NULL,
  PRIMARY KEY  (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `course` (
  `id` INT(10) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) DEFAULT NULL,
  PRIMARY KEY  (`id`)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;