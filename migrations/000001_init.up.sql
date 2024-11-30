CREATE TABLE payment_requests (
    id SERIAL PRIMARY KEY,                          
    request_id INTEGER NOT NULL, 
    amount INT NOT NULL,       
    payment_status VARCHAR(20),
    callback_url TEXT,                 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_id) REFERENCES "rent-request-service".rent_requests(id) ON DELETE SET NULL
);