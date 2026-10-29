CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE authors (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    bio TEXT
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL CHECK (role IN ('reader', 'librarian', 'admin'))
);

CREATE TABLE books (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    isbn TEXT NOT NULL UNIQUE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT
);

CREATE INDEX idx_books_category_id ON books(category_id);

CREATE TABLE book_authors (
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES authors(id) ON DELETE RESTRICT,
    PRIMARY KEY (book_id, author_id)
);

CREATE INDEX idx_book_authors_author_id ON book_authors(author_id);

CREATE TABLE copies (
    id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    inventory_no TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'borrowed', 'lost'))
);

CREATE INDEX idx_copies_book_id ON copies(book_id);
