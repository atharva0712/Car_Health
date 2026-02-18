CREATE DATABASE IF NOT EXISTS car_monitor;
USE car_monitor;

CREATE TABLE cars (
    vin VARCHAR(50) PRIMARY KEY,
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE attribute_status (
    vin VARCHAR(50) PRIMARY KEY,
    tire_pressure_value DECIMAL(5,2) NOT NULL,
    tire_pressure_status ENUM('NORMAL','WARNING','CRITICAL') NOT NULL,
    engine_temp_value DECIMAL(5,2) NOT NULL,
    engine_temp_status ENUM('NORMAL','WARNING','CRITICAL') NOT NULL,
    battery_voltage_value DECIMAL(5,2) NOT NULL,
    battery_voltage_status ENUM('NORMAL','WARNING','CRITICAL') NOT NULL,
    oil_level_value DECIMAL(5,2) NOT NULL,
    oil_level_status ENUM('NORMAL','WARNING','CRITICAL') NOT NULL,
    brake_condition_value DECIMAL(5,2) NOT NULL,
    brake_condition_status ENUM('NORMAL','WARNING','CRITICAL') NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (vin) REFERENCES cars(vin) ON DELETE CASCADE
);

CREATE TABLE health_events (
    id INT PRIMARY KEY AUTO_INCREMENT,
    vin VARCHAR(50) NOT NULL,
    attribute_name VARCHAR(50) NOT NULL,
    event_type ENUM('WARNING','CRITICAL') NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NULL,
    duration_seconds INT NULL,
    FOREIGN KEY (vin) REFERENCES cars(vin) ON DELETE CASCADE
);