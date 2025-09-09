-- Enable UUID extension (if not already enabled)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fullname TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Hotels Table
CREATE TABLE IF NOT EXISTS hotels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    manager_id UUID NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_hotel_manager FOREIGN KEY (manager_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Rooms Table
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id UUID NOT NULL,
    available_quantity INT NOT NULL CHECK (available_quantity >= 0),
    room_category TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_room_hotel FOREIGN KEY (hotel_id)
        REFERENCES hotels(id)
        ON DELETE CASCADE
);

-- Bookings Table
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    hotel_id UUID NOT NULL,
    checkin TIMESTAMPTZ NOT NULL,
    checkout TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_booking_user FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_booking_hotel FOREIGN KEY (hotel_id)
        REFERENCES hotels(id)
        ON DELETE CASCADE
);

-- BookedRooms Table
CREATE TABLE IF NOT EXISTS booked_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL,
    room_type TEXT NOT NULL,
    room_quantity INT NOT NULL CHECK (room_quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_booked_room_booking FOREIGN KEY (booking_id)
        REFERENCES bookings(id)
        ON DELETE CASCADE
);
