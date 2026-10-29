INSERT INTO categories (name) VALUES
    ('Programming'),
    ('Fiction'),
    ('History');

INSERT INTO authors (full_name, bio) VALUES
    ('Alan Turing', 'Computing pioneer'),
    ('Ada Lovelace', 'First programmer'),
    ('Ursula K. Le Guin', 'Novelist'),
    ('Mary Beard', 'Historian'),
    ('Donald Knuth', 'Algorithms');

INSERT INTO books (title, isbn, category_id) VALUES
    ('The Annotated Turing', '9780470229057', 1),
    ('Notes on Programming', '9781111111111', 1),
    ('The Left Hand of Darkness', '9780441478125', 2),
    ('A Wizard of Earthsea', '9780553383041', 2),
    ('SPQR', '9781631492228', 3),
    ('Women and Power', '9781631494758', 3),
    ('The Art of Computer Programming', '9780201896831', 1),
    ('Concrete Mathematics', '9780201558029', 1);

INSERT INTO book_authors (book_id, author_id) VALUES
    (1, 1),
    (1, 2),
    (2, 1),
    (3, 3),
    (4, 3),
    (5, 4),
    (6, 4),
    (7, 5),
    (8, 5);

INSERT INTO copies (book_id, inventory_no, status) VALUES
    (1, 'INV-001', 'available'),
    (1, 'INV-002', 'available'),
    (1, 'INV-003', 'lost'),
    (2, 'INV-004', 'available'),
    (2, 'INV-005', 'borrowed'),
    (2, 'INV-006', 'available'),
    (3, 'INV-007', 'available'),
    (3, 'INV-008', 'available'),
    (4, 'INV-009', 'available'),
    (4, 'INV-010', 'lost'),
    (5, 'INV-011', 'available'),
    (5, 'INV-012', 'available'),
    (5, 'INV-013', 'borrowed'),
    (6, 'INV-014', 'available'),
    (6, 'INV-015', 'available'),
    (7, 'INV-016', 'available'),
    (7, 'INV-017', 'available'),
    (7, 'INV-018', 'available'),
    (8, 'INV-019', 'available'),
    (8, 'INV-020', 'borrowed'),
    (8, 'INV-021', 'available');
