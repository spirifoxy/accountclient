Feature: check health status
    In order to use the API client
    As a developer using F3 API
    I need to be sure API is alive

    Scenario: should get success response
        When I call IsHealthy method
        Then I should receive true as response
