DROP TABLE IF EXISTS `chat`;
CREATE TABLE `chat` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `time` varchar(30) NOT NULL,
  `user` varchar(30) NOT NULL,
  `message` varchar(30) NOT NULL,
  `hash` varchar(255) NOT NULL,
  `proof` varchar(255) NOT NULL,
  `blockchain` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;