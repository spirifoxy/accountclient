Feature: create an account
    In order to use account API
    As a developer using F3 API
    I need to be able to create new account entities

    Scenario: should create an account
        When I try to create the account with uuid 11111111-1111-1111-1111-111111111111
        Then the response is success
        And I get the account

    Scenario: should create several different accounts
        When I try to create the account with uuid 11111111-1111-1111-1111-111111111111
        And I try to create the account with uuid 22222222-2222-2222-2222-222222222222
        Then the response is success
        And I get the account

    Scenario: should not create duplicated account
        When I try to create the account with uuid 11111111-1111-1111-1111-111111111111
        And I try to create the account with uuid 11111111-1111-1111-1111-111111111111
        Then the response code should be 409
        And I get the request status error

    Scenario: should not create corrupted account with no name
        When I accidentally try to create an account with no name
        Then the response code should be 400
        And I get the request status error

    Scenario: should not create corrupted account with broken country
        When I accidentally try to create an account with broken country
        Then the response code should be 400
        And I get the request status error

    Scenario: should not create corrupted account with organization id
        When I accidentally try to create an account with broken organization id
        Then the response code should be 400
        And I get the request status error
