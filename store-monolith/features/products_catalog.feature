Feature: Products Catalog
    The service should be able to
    list all products, as well as
    search and get details of the
    products

  Scenario Outline: List all products
    Given products:
      | id | name      | description           | price | currency |
      |  1 | Product 1 | Product description 1 | 19.99 | EUR      |
      |  2 | Product 2 | Product description 2 | 34.99 | HUF      |
    When http get "/products"
    Then the response should be 200
    And <field> should be <value>

    Examples: 
      | field       | value     |
      | [id=1].name | Product 1 |
