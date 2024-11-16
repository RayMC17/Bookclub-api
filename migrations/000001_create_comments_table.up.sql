CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    authors TEXT[] NOT NULL,
    isbn VARCHAR(20) UNIQUE NOT NULL,
    publication_date DATE NOT NULL,
    genre VARCHAR(50),
    description TEXT,
    average_rating FLOAT DEFAULT 0
);


CREATE TABLE reading_lists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_by INT REFERENCES users(id) ON DELETE CASCADE,
    books INT[] DEFAULT '{}',
    status VARCHAR(20) CHECK (status IN ('currently reading', 'completed')) NOT NULL
);

CREATE TABLE book_reviews (
    id SERIAL PRIMARY KEY,
    book_id INT REFERENCES books(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    review_date DATE DEFAULT CURRENT_DATE
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    reading_lists INT[] DEFAULT '{}',
    reviews INT[] DEFAULT '{}'
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author VARCHAR(255) NOT NULL,
    rating DOUBLE PRECISION NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
