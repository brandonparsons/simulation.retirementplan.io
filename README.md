Retirement Simulation - Golang
===============================

Runtime Dependencies
--------------------

For zero-downtime re-deploy etc:
- `einhorn` (Ruby)

Running
-------

- `gem install einhorn`
- `einhorn -b 127.0.0.1:1234 -m manual go_retirement_simulation`

Development
------------

- `get get .`
- `go run main.go` or `gin`

Examples
--------

```
http POST 'http://localhost:3000/simulation' \
'Authorization: abcd' \
in_todays_dollars:=true \
number_of_trials:=2 \
selected_portfolio_weights:='{ "INTL-BOND": 0.65, "US-REALESTATE": 0.30, "CDN-REALESTATE": 0.05 }' \ 
inflation:='{"mean": 0.00046346514957523,"std_dev": 0.00024792742828969}' \
real_estate:='{"mean": 0.0029064094738571,"std_dev": 0.014660011854061}' \

simulation_parameters:='{ "male": true, "married": true, "retired": false, "male_age": 29, "retirement_age_male": 62, "female_age": 30, "retirement_age_female": 35, "expenses_multiplier": 1.6, "fraction_single_income": 65, "starting_assets": 125000, "income": 120000, "current_tax": 35, "salary_increase": 3, "income_inflation_index": 20, "expenses_inflation_index": 100, "retirement_income": 12000, "retirement_expenses": 80, "retirement_tax": 25, "life_insurance": 250000, "include_home": true, "home_value": 550000, "sell_house_in": 25, "new_home_relative_value": 65 }' \
expenses:='[{"amount": 300,"frequency": "monthly","onetime_on": null,"ends": null},{"amount": 25000,"frequency": "onetime","onetime_on": 1461564000000,"ends": null}]' \
```


```ruby
require 'faraday'
require 'json'

payload = {
    in_todays_dollars: true,
    number_of_trials: 2,
    selected_portfolio_weights: { 
        "INTL-BOND" => 0.65, 
        "US-REALESTATE" => 0.30, 
        "CDN-REALESTATE" => 0.05
    },
    asset_performance_data: {
        "INTL-BOND" =>  {
            mean:    0.0003,
            std_dev: 0.0002,
        }, 
        "US-REALESTATE" =>  {
            mean:    0.0004,
            std_dev: 0.00025,
        }, 
        "CDN-REALESTATE" =>  {
            mean:    0.0005,
            std_dev: 0.00021,
        }, 
    },
    cholesky_decomposition: [
        0.0094794922, 
        0.0, 
        0.0, 
        -7.36e-05, 
        0.0055677999, 
        0.0, 
        0.0050681903, 
        -0.0004821709, 
        0.013367741
    ],
    inflation: {
        mean: 0.00046346514957523,
        std_dev: 0.00024792742828969
    },
    real_estate: {
        mean: 0.0029064094738571,
        std_dev: 0.014660011854061
    },
    expenses: [
        {amount: 300, frequency: 'monthly', onetime_on: nil, ends: nil},
        {amount: 25000, frequency: 'onetime', onetime_on: 1461564000000, ends: nil}
    ], 
    simulation_parameters: {
        male: true,
        married: true,
        retired: false,
        male_age: 29,
        retirement_age_male: 62,
        female_age: 30,
        retirement_age_female: 35,
        expenses_multiplier: 1.6,
        fraction_single_income: 65,
        starting_assets: 125000,
        income: 120000,
        current_tax: 35,
        salary_increase: 3,
        income_inflation_index: 20,
        expenses_inflation_index: 100,
        retirement_income: 12000,
        retirement_expenses: 80,
        retirement_tax: 25,
        life_insurance: 250000,
        include_home: true,
        home_value: 550000,
        sell_house_in: 25,
        new_home_relative_value: 65 
    }
}

conn = Faraday.new(:url => 'http://localhost:3000') do |faraday|
  # faraday.response :logger                # log requests to STDOUT
  faraday.adapter  Faraday.default_adapter  # make requests with Net::HTTP
end

def get_response(conn, payload)
    response = JSON.parse(conn.post do |req|
      req.url '/simulation'
      req.headers['Authorization'] = 'abcd'
      req.headers['Content-Type'] = 'application/json'
      req.body = payload.to_json
    end.body)
    return response
end
```
