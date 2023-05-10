-- MySQL dump 10.13  Distrib 5.7.26, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: vita_task
-- ------------------------------------------------------
-- Server version	5.7.26

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES UTF8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `vita_task`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `vita_task` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `vita_task`;

--
-- Table structure for table `vt_dialog`
--

DROP TABLE IF EXISTS `vt_dialog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_dialog` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `create_time` bigint(20) unsigned DEFAULT '0',
  `update_time` bigint(20) unsigned DEFAULT '0',
  `type` varchar(30) DEFAULT '',
  `name` longtext,
  `deleted_at` datetime(3) DEFAULT NULL,
  `last_at` bigint(20) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_dialog`
--

LOCK TABLES `vt_dialog` WRITE;
/*!40000 ALTER TABLE `vt_dialog` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_dialog` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_dialog_msg`
--

DROP TABLE IF EXISTS `vt_dialog_msg`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_dialog_msg` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `create_time` bigint(20) unsigned DEFAULT '0',
  `update_time` bigint(20) unsigned DEFAULT '0',
  `dialog_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `type` varchar(30) DEFAULT '',
  `content` longtext,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_vt_dialog_msg_dialog` (`dialog_id`),
  CONSTRAINT `fk_vt_dialog_msg_dialog` FOREIGN KEY (`dialog_id`) REFERENCES `vt_dialog` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_dialog_msg`
--

LOCK TABLES `vt_dialog_msg` WRITE;
/*!40000 ALTER TABLE `vt_dialog_msg` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_dialog_msg` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_dialog_user`
--

DROP TABLE IF EXISTS `vt_dialog_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_dialog_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `create_time` bigint(20) unsigned DEFAULT '0',
  `update_time` bigint(20) unsigned DEFAULT '0',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `dialog_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_dialog_user`
--

LOCK TABLES `vt_dialog_user` WRITE;
/*!40000 ALTER TABLE `vt_dialog_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_dialog_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_org_user`
--

DROP TABLE IF EXISTS `vt_org_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_org_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) DEFAULT '0',
  `org_id` int(11) DEFAULT '0',
  `role` int(11) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`,`org_id`,`role`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_org_user`
--

LOCK TABLES `vt_org_user` WRITE;
/*!40000 ALTER TABLE `vt_org_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_org_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_organization`
--

DROP TABLE IF EXISTS `vt_organization`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_organization` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `type` tinyint(4) DEFAULT '0' COMMENT '分类。1-公司 2-部门 3-团队',
  `addr` varchar(255) NOT NULL DEFAULT '',
  `register_date` datetime DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `deleted` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `parent_id` (`parent_id`,`type`),
  KEY `deleted` (`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组织架构表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_organization`
--

LOCK TABLES `vt_organization` WRITE;
/*!40000 ALTER TABLE `vt_organization` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_organization` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_project`
--

DROP TABLE IF EXISTS `vt_project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_project` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '项目名',
  `create_time` bigint(20) DEFAULT '0',
  `update_time` bigint(20) DEFAULT '0',
  `deleted_at` datetime DEFAULT NULL,
  `complete` int(11) NOT NULL DEFAULT '0' COMMENT '已完成任务数量',
  `archive` tinyint(4) DEFAULT '0' COMMENT '归档',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_project`
--

LOCK TABLES `vt_project` WRITE;
/*!40000 ALTER TABLE `vt_project` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_project` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_project_member`
--

DROP TABLE IF EXISTS `vt_project_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_project_member` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '项目ID',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `role` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '角色',
  PRIMARY KEY (`id`),
  KEY `project_id` (`project_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_project_member`
--

LOCK TABLES `vt_project_member` WRITE;
/*!40000 ALTER TABLE `vt_project_member` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_project_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_task`
--

DROP TABLE IF EXISTS `vt_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_task` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '项目ID',
  `group_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '项目任务组ID',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '任务标题',
  `describe` longtext COMMENT '任务描述',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '任务状态',
  `level` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '任务紧急度',
  `complete_date` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '完成时间',
  `archived_date` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '归档时间',
  `start_date` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '计划开始时间',
  `end_date` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '计划结束时间',
  `enclosure_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '附件数量',
  `dialog_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '对话ID',
  `create_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `project_id` (`project_id`,`group_id`,`status`,`level`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_task`
--

LOCK TABLES `vt_task` WRITE;
/*!40000 ALTER TABLE `vt_task` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_task_files`
--

DROP TABLE IF EXISTS `vt_task_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_task_files` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '项目ID',
  `task_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `filename` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT '路径',
  `md5` varchar(50) NOT NULL DEFAULT '' COMMENT '文件MD5值',
  `size` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '文件大小(B)',
  `ext` varchar(255) NOT NULL DEFAULT '' COMMENT '文件扩展名',
  `download` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '下载次数',
  `thumb` varchar(255) NOT NULL DEFAULT '' COMMENT '缩略图',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `project_id` (`project_id`,`task_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_task_files`
--

LOCK TABLES `vt_task_files` WRITE;
/*!40000 ALTER TABLE `vt_task_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_task_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_task_group`
--

DROP TABLE IF EXISTS `vt_task_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_task_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '项目ID',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '任务组名称',
  `create_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_task_group`
--

LOCK TABLES `vt_task_group` WRITE;
/*!40000 ALTER TABLE `vt_task_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_task_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_task_log`
--

DROP TABLE IF EXISTS `vt_task_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_task_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(10) unsigned NOT NULL DEFAULT '0',
  `operate_type` char(30) NOT NULL DEFAULT '' COMMENT '操作类型',
  `operator` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '操作人员ID',
  `operate_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '操作时间',
  `message` longtext COMMENT '日志信息',
  `create_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`) USING BTREE
) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_task_log`
--

LOCK TABLES `vt_task_log` WRITE;
/*!40000 ALTER TABLE `vt_task_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_task_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_task_member`
--

DROP TABLE IF EXISTS `vt_task_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_task_member` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `role` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '角色',
  PRIMARY KEY (`id`),
  KEY `task_id` (`task_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_task_member`
--

LOCK TABLES `vt_task_member` WRITE;
/*!40000 ALTER TABLE `vt_task_member` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_task_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user`
--

DROP TABLE IF EXISTS `vt_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vt_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_login` varchar(60) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pass` varchar(64) NOT NULL DEFAULT '' COMMENT '登录密码;',
  `user_status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '用户状态;1:正常,2:禁用',
  `user_nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `sex` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别;0:保密,1:男,2:女',
  `birthday` datetime DEFAULT NULL COMMENT '生日',
  `user_email` varchar(100) NOT NULL DEFAULT '' COMMENT '用户登录邮箱',
  `avatar` varchar(1024) NOT NULL DEFAULT '' COMMENT '用户头像',
  `signature` varchar(255) NOT NULL DEFAULT '' COMMENT '个性签名',
  `user_activation_key` varchar(60) NOT NULL DEFAULT '' COMMENT '激活码',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '用户手机号',
  `lock_time` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '登陆错误锁定结束时间',
  `error_sum` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '登陆错误次数',
  `first` tinyint(3) unsigned DEFAULT '1' COMMENT '是否首次登录系统',
  `last_edit_pass` bigint(20) DEFAULT '0' COMMENT '最后一次修改密码的时间',
  `last_login_ip` varchar(15) NOT NULL DEFAULT '' COMMENT '最后登录ip',
  `last_login_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `openid` varchar(64) DEFAULT '' COMMENT '微信openid',
  `super` tinyint(1) unsigned DEFAULT '0' COMMENT '超级管理员',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '注册时间',
  `update_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '信息更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_login` (`user_login`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user`
--

LOCK TABLES `vt_user` WRITE;
/*!40000 ALTER TABLE `vt_user` DISABLE KEYS */;
INSERT INTO `vt_user` VALUES (1,'admin','88485f172d58f133b1f611b411ee646e',1,'超级管理员',1,'2022-02-23 00:00:00','45451212@qq.com','/uploads\\20230415\\0305741cd93fe108590b36bfe2feadd4.jpg','只因你太美','','15889891212',0,0,1,1681715241764,'',0,'',1,1681723693,1681723693);
/*!40000 ALTER TABLE `vt_user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-05-10 15:41:35
