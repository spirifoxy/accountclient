Feature: fetch an account
    In order to use account API
    As a developer using F3 API
    I need to be able to query specific account entity

    Scenario: should get an account
        Given the account exists with uuid 11111111-1111-1111-1111-111111111111
        When I try to fetch the account with uuid 11111111-1111-1111-1111-111111111111
        Then the response is success
        And I get the account

    Scenario: should not get missing account
        When I try to fetch the account with uuid 11111111-1111-1111-1111-111111111111
        Then the response code should be 404
        And I get the request status error

    Scenario: should get recently created account
        Given I created an account with uuid 11111111-1111-1111-1111-111111111111
        When I try to fetch the account with uuid 11111111-1111-1111-1111-111111111111
        Then the response is success
        And I get the account

    Scenario: should not get removed account
        Given I created an account with uuid 11111111-1111-1111-1111-111111111111
        And I removed the account with uuid 11111111-1111-1111-1111-111111111111
        When I try to fetch the account with uuid 11111111-1111-1111-1111-111111111111
        Then the response code should be 404
        And I get the request status error
