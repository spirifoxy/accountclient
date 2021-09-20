Feature: delete an account
    In order to use account API
    As a developer using F3 API
    I need to be able to remove specific account entity

    Scenario: should remove an account
        Given the account exists with uuid 11111111-1111-1111-1111-111111111111
        When I try to remove an account with uuid 11111111-1111-1111-1111-111111111111
        Then the response is success

    Scenario: should not remove unexistent account
        When I try to remove an account with uuid 11111111-1111-1111-1111-111111111111
        Then the response code should be 404
        And I get the request status error

    Scenario: should not remove account of wrong version
        Given the account exists with uuid 11111111-1111-1111-1111-111111111111
        When I try to remove the modified account with uuid 11111111-1111-1111-1111-111111111111
        Then the response code should be 409
        And I get the request status error
