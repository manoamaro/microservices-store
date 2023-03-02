DELETE FROM auths_domains
WHERE (auth_id, domain_id) IN (
    SELECT a.id, d.id FROM auths a, domains d
    WHERE a.email = 'admin@admin.com'
      AND d.domain in ('auth_admin', 'products_admin', 'orders_admin', 'cart_admin')
);

DELETE FROM domains
WHERE "domain" IN ('auth_admin', 'products_admin', 'orders_admin', 'cart_admin');

DELETE FROM auths WHERE email = 'admin@admin.com';