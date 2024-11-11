-- MySQL dump 10.13  Distrib 8.0.38, for Win64 (x86_64)
--
-- Host: localhost    Database: users-api
-- ------------------------------------------------------
-- Server version	8.0.39

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users`
(
    `user_id`  int          NOT NULL AUTO_INCREMENT,
    `email`    varchar(100) NOT NULL,
    `password` varchar(100) NOT NULL,
    `nombre`   varchar(100) NOT NULL,
    `apellido` varchar(100) NOT NULL,
    `admin`    tinyint(1) NOT NULL,
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `uni_users_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK
TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users`
VALUES (1, 'sofiaolivetoo@gmail.com', 'ba77a5448b1208afe6effd5194c2a8b6', 'Sofia', 'Oliveto', 0),
       (2, 'juanlopez@gmail.com', 'f5737d25829e95b9c234b7fa06af8736', 'Juan', 'Lopez', 1),
       (3, 'constanzastrumia@gmail.com', 'febf04180a62e8710868cafd8741515f', 'Constanza', 'Strumia', 0),
       (4, 'margarita@gmail.com', '828fca74e9e1d7e55b76d46a304b5f55', 'Margarita', 'de Marcos', 1),
       (5, 'pedro@gmail.com', 'd3ce9efea6244baa7bf718f12dd0c331', 'Pedro', 'Juarez', 0),
       (6, 'josefinagonzalez@gmail.com', 'e577bc7b26b52afe6a33f02513b86b5c', 'Josefina', 'Gonzalez', 0),
       (7, 'ramiropaez@gmail.com', '8e7a60d71791c1febdbc4998c963e87e', 'Ramiro', 'Paez', 0),
       (8, 'gustavojacobo@gmail.com', '0805446e686aa72d45f9583f2d6cedef', 'Gustavo', 'Jacobo', 1),
       (9, 'matigarcia@gmail.com', '0596f701227172915b2862b95b4c2e1a', 'Matias', 'Portillo', 0),
       (10, 'juliomansilla@gmail.com', '16880e98af692b72ce3ba695654ee306', 'Julio', 'Mansilla', 0),
       (11, 'santiportillo@gmail.com', 'd5116c2a9607b0ea07d425506f610467', 'Santiago', 'Portillo', 0),
       (12, 'nicolasfigueroa@gmail.com', '305735d035e7f7381d64d179126ff6d9', 'Nicolas', 'Figueroa', 0),
       (13, 'agostinacisneros@gmail.com', '0a9827114b460ae8f2f96c5e8893c90d', 'Agostina', 'Cisneros', 0),
       (14, 'luciabernardi@gmail.com', 'b9ce57d6d6c2d6fda25d80da5a00a7d1', 'Lucia', 'Bernardi', 0),
       (15, 'pauladominguez@gmail.com', 'ca46a79286419c05172ca7b010a59d3c', 'Paula', 'Dominguez', 0),
       (16, 'luciovelarde@gmail.com', 'f4a1f4d408e436f3c294bf2ae346b3d7', 'Lucio', 'Velarde', 0),
       (17, 'rodolfoperez@gmail.com', '4a2b3910b547e5212914378adaf76aac', 'Rodolfo', 'Perez', 1),
       (18, 'sebastiancolidio@gmail.com', '5a7c2cf0d17f9d32c87de8efb8e689d6', 'Sebastian', 'Colidio', 1),
       (19, 'lucasbeltran@gmail.com', '6d16ba70238c92a03ac04c7c86eb79e7', 'Lucas', 'Beltran', 1),
       (20, 'chilenodiaz@gmail.com', '4494d10dc9752cba4083ce2cf8983d2c', 'Paulo', 'Diaz', 0),
       (21, '', '54afa97407b6b1aac35e8550365cc7c8', '', '', 0);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK
TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-11-11 12:14:25
