/*M!999999\- enable the sandbox mode */ 

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
DROP TABLE IF EXISTS `pws_latest`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pws_latest` (
  `station_ref_id` bigint(20) unsigned NOT NULL,
  `temp_c` decimal(5,1) DEFAULT NULL,
  `humidity` decimal(5,1) DEFAULT NULL,
  `obs_time_utc` datetime DEFAULT NULL,
  `fetched_at_utc` datetime NOT NULL,
  `stale` tinyint(1) NOT NULL DEFAULT 0,
  `invalid` tinyint(1) NOT NULL DEFAULT 0,
  `error_message` varchar(255) DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`station_ref_id`),
  KEY `pws_latest_fetched_at_index` (`fetched_at_utc`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `pws_stations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pws_stations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `station_id` varchar(64) NOT NULL,
  `name` varchar(120) NOT NULL,
  `lat` decimal(9,6) DEFAULT NULL,
  `lon` decimal(9,6) DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `display_order` int(11) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `pws_stations_station_id_unique` (`station_id`),
  KEY `pws_stations_active_order_index` (`is_active`,`display_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `weather`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `weather` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `measured_at` datetime NOT NULL COMMENT 'Time when the weather data was measured',
  `temperature` decimal(5,1) NOT NULL COMMENT 'Temperature in °C',
  `pressure` decimal(7,1) NOT NULL COMMENT 'Pressure in hPa',
  `humidity` decimal(5,1) NOT NULL COMMENT 'Humidity in %',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `weather_daily`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `weather_daily` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL COMMENT 'Date of the measurement',
  `sea_temperature` decimal(5,1) DEFAULT NULL COMMENT 'Sea water temperature in °C (manually measured)',
  `avg_temperature` decimal(5,1) DEFAULT NULL,
  `min_temperature` decimal(5,1) DEFAULT NULL,
  `max_temperature` decimal(5,1) DEFAULT NULL,
  `avg_pressure` decimal(7,1) DEFAULT NULL,
  `min_pressure` decimal(7,1) DEFAULT NULL,
  `max_pressure` decimal(7,1) DEFAULT NULL,
  `avg_humidity` decimal(5,1) DEFAULT NULL,
  `min_humidity` decimal(5,1) DEFAULT NULL,
  `max_humidity` decimal(5,1) DEFAULT NULL,
  `samples_count` int(11) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `weather_daily_date_unique` (`date`),
  KEY `weather_daily_date_index` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `weather_hourly`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `weather_hourly` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL COMMENT 'Date of the measurement',
  `hour` tinyint(4) NOT NULL COMMENT 'Hour (0-23)',
  `avg_temperature` decimal(5,1) NOT NULL COMMENT 'Average temperature in °C',
  `avg_pressure` decimal(7,1) NOT NULL COMMENT 'Average pressure in hPa',
  `avg_humidity` decimal(5,1) NOT NULL COMMENT 'Average humidity in %',
  `samples_count` int(11) NOT NULL COMMENT 'Number of measurements used for average',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `weather_hourly_date_hour_unique` (`date`,`hour`),
  KEY `weather_hourly_date_index` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `weather_monthly`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `weather_monthly` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `year` int(11) NOT NULL COMMENT 'Year (e.g., 2025)',
  `month` tinyint(4) NOT NULL COMMENT 'Month (1-12)',
  `avg_temperature` decimal(5,1) NOT NULL COMMENT 'Average temperature in °C',
  `min_temperature` decimal(5,1) NOT NULL COMMENT 'Minimum temperature in °C',
  `max_temperature` decimal(5,1) NOT NULL COMMENT 'Maximum temperature in °C',
  `avg_pressure` decimal(7,1) NOT NULL COMMENT 'Average pressure in hPa',
  `min_pressure` decimal(7,1) NOT NULL COMMENT 'Minimum pressure in hPa',
  `max_pressure` decimal(7,1) NOT NULL COMMENT 'Maximum pressure in hPa',
  `avg_humidity` decimal(5,1) NOT NULL COMMENT 'Average humidity in %',
  `min_humidity` decimal(5,1) NOT NULL COMMENT 'Minimum humidity in %',
  `max_humidity` decimal(5,1) NOT NULL COMMENT 'Maximum humidity in %',
  `samples_count` int(11) NOT NULL COMMENT 'Number of measurements used',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `weather_monthly_year_month_unique` (`year`,`month`),
  KEY `weather_monthly_year_month_index` (`year`,`month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `weather_weekly`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `weather_weekly` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `year` int(11) NOT NULL COMMENT 'Year (e.g., 2025)',
  `week` tinyint(4) NOT NULL COMMENT 'ISO week number (1-53)',
  `week_start` date NOT NULL COMMENT 'Monday of the week',
  `week_end` date NOT NULL COMMENT 'Sunday of the week',
  `avg_temperature` decimal(5,1) NOT NULL COMMENT 'Average temperature in °C',
  `min_temperature` decimal(5,1) NOT NULL COMMENT 'Minimum temperature in °C',
  `max_temperature` decimal(5,1) NOT NULL COMMENT 'Maximum temperature in °C',
  `avg_pressure` decimal(7,1) NOT NULL COMMENT 'Average pressure in hPa',
  `min_pressure` decimal(7,1) NOT NULL COMMENT 'Minimum pressure in hPa',
  `max_pressure` decimal(7,1) NOT NULL COMMENT 'Maximum pressure in hPa',
  `avg_humidity` decimal(5,1) NOT NULL COMMENT 'Average humidity in %',
  `min_humidity` decimal(5,1) NOT NULL COMMENT 'Minimum humidity in %',
  `max_humidity` decimal(5,1) NOT NULL COMMENT 'Maximum humidity in %',
  `samples_count` int(11) NOT NULL COMMENT 'Number of measurements used',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `weather_weekly_year_week_unique` (`year`,`week`),
  KEY `weather_weekly_year_week_index` (`year`,`week`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!50001 ALTER TABLE `pws_latest`
  ADD CONSTRAINT `pws_latest_station_ref_id_fk` FOREIGN KEY (`station_ref_id`) REFERENCES `pws_stations` (`id`) ON DELETE CASCADE */;
DROP TABLE IF EXISTS `tide_events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `tide_events` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `date_local` date NOT NULL COMMENT 'Local date (Europe/Madrid)',
  `location_key` varchar(64) NOT NULL,
  `event_type` varchar(8) NOT NULL COMMENT 'HIGH|LOW',
  `event_time_local` datetime NOT NULL COMMENT 'Local datetime (Europe/Madrid)',
  `height_m` decimal(8,3) NOT NULL,
  `source` varchar(32) NOT NULL COMMENT 'puertos|open_meteo',
  `confidence` tinyint(3) unsigned NOT NULL,
  `fetched_at` datetime NOT NULL COMMENT 'UTC fetch timestamp',
  `raw_json` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `tide_events_unique` (`date_local`,`location_key`,`event_type`,`event_time_local`),
  KEY `tide_events_lookup_idx` (`date_local`,`location_key`,`source`,`fetched_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
