INSERT INTO domains (domain) VALUES ('auth_admin'),
                                    ('products_admin'),
                                    ('orders_admin'),
                                    ('cart_admin'),
                                    ('inventory');

INSERT INTO auths (email, password, salt)
VALUES ('admin@admin.com','d82494f05d6917ba02f7aaa29689ccb444bb73f20380876cb05d1f37537b7892', 'admin');

INSERT INTO audiences(auth_id, domain_id) (
    select a.id, d.id from auths a, domains d
    where a.email = 'admin@admin.com'
      and d.domain in (select "domain" from domains)
);
