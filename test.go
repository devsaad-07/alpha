package main

func DefineRuleSet() string {
	return `
    rule "Discount for loyal customers"
        when
            $customer: Customer(orderCount >= 5)
        then
            $customer.SetDiscount(0.1);
    end
    `
}
