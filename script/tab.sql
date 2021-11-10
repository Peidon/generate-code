CREATE TABLE `user_field_extend_tab_00000000` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `saas_id`        bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'saas id',
  `region`             varchar(10) NOT NULL DEFAULT '' COMMENT 'region',
  `cs_user_id` bigint(20) NOT NULL COMMENT 'cs user id',
  `field_id` bigint(20) NOT NULL COMMENT 'user define field id',
  `field_value` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'field value',
  `status_flag` tinyint(4) NOT NULL COMMENT 'status flag: 1-normal -1-deleted',
  `mtime` int(10) unsigned NOT NULL COMMENT 'modify time',
  `ctime` int(10) unsigned NOT NULL COMMENT 'create time',
  PRIMARY KEY (`id`),
  KEY `idx_cs_user_id_field_id` (`saas_id`, `region`,`cs_user_id`,`field_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4