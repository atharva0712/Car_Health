CREATE DATABASE IF NOT EXISTS car_monitor;
USE car_monitor;

CREATE TABLE tenants (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    api_key_hash CHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cars (
    tenant_id BIGINT NOT NULL,
    vin VARCHAR(50) NOT NULL,
    simulation_enabled BOOLEAN DEFAULT TRUE,
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tenant_id, vin),
    CONSTRAINT fk_cars_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE TABLE IF NOT EXISTS attribute_status (
    tenant_id BIGINT NOT NULL,
    vin VARCHAR(50) NOT NULL,
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
    PRIMARY KEY (tenant_id, vin),
    CONSTRAINT fk_attribute_status_car FOREIGN KEY (tenant_id, vin)
        REFERENCES cars(tenant_id, vin) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS health_events (
    tenant_id BIGINT NOT NULL,
    vin VARCHAR(50) NOT NULL,
    attribute_name VARCHAR(50) NOT NULL,
    event_type ENUM('WARNING','CRITICAL') NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NULL,
    duration_seconds INT NULL,
    PRIMARY KEY (tenant_id, vin, attribute_name, start_time),
    CONSTRAINT fk_health_events_car FOREIGN KEY (tenant_id, vin)
        REFERENCES cars(tenant_id, vin) ON DELETE CASCADE
);