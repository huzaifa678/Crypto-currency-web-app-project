package val

import (
	"fmt"
	"net/mail"
	"regexp"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/shopspring/decimal"
)

var (
	isValidUsername   = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	//isValidFullName   = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
	isValidCurrency   = regexp.MustCompile(`^[A-Z]{3}$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}

func ValidateUserRole(role pb.UserRole) error {
	switch role {
	case pb.UserRole_USER_ROLE_ADMIN, pb.UserRole_USER_ROLE_USER:
		return nil
	default:
		return fmt.Errorf("invalid user role")
	}
}

func ValidateOrderType(orderType db.OrderType) error {
	switch orderType {
	case db.OrderTypeBuy, db.OrderTypeSell:
		return nil
	default:
		return fmt.Errorf("invalid order type")
	}
}

func ValidateCurrency(value string) error {
	if !isValidCurrency(value) {
		return fmt.Errorf("must be a valid 3-letter currency code")
	}
	return nil
}

func ValidateMarket(marketID, baseCurrency, quoteCurrency string, minOrderAmount decimal.Decimal, pricePrecision int32) error {
	if err := ValidateString(marketID, 1, 50); err != nil {
		return fmt.Errorf("market_id %v", err)
	}
	if err := ValidateCurrency(baseCurrency); err != nil {
		return fmt.Errorf("base_currency %v", err)
	}
	if err := ValidateCurrency(quoteCurrency); err != nil {
		return fmt.Errorf("quote_currency %v", err)
	}
	if err := validateDecimal(minOrderAmount); err != nil {
		return fmt.Errorf("min_order_amount %v", err)
	}
	if pricePrecision < 0 {
		return fmt.Errorf("price_precision must be a non-negative integer")
	}
	return nil
}

func ValidateLoginUserRequest(email, password string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidateString(password, 6, 100); err != nil {
		return err
	}
	return nil
}

func ValidateCreateUserRequest(username string, email string, password string, role pb.UserRole) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidateString(password, 6, 100); err != nil {
		return err
	}
	if err := ValidateUserRole(role); err != nil {
		return err
	}
	return nil
}

func ValidateGetRequest(id string) error {
	return ValidateString(id, 1, 50)
}

func ValidateDeleteRequest(id string) error {
	return ValidateString(id, 1, 50)
}

func ValidateCreateMarketRequest(baseCurrency, quoteCurrency string, minOrderAmount decimal.Decimal, pricePrecision int32) error {
	return ValidateMarket("dummy", baseCurrency, quoteCurrency, minOrderAmount, pricePrecision)
}

func ValidateCreateOrderRequest(userEmail string, marketID string, price decimal.Decimal, amount decimal.Decimal, orderType pb.OrderType) error {
	if err := ValidateEmail(userEmail); err != nil {
		return err
	}

	if err := ValidateString(marketID, 1, 50); err != nil {
		return err
	}

	if err := validateDecimal(price); err != nil {
		return err
	}

	if err := validateDecimal(amount); err != nil {
		return err
	}

	if orderType < 0 || orderType > 1 {
		return fmt.Errorf("order_type must be 0 (BUY) or 1 (SELL)")
	}
	return nil
}


func ValidateUpdateUser(userID, password string) error {
	if err := ValidateString(userID, 1, 50); err != nil {
		return err
	}
	return ValidateString(password, 6, 100)
}

func ValidateUpdateWalletRequest(walletID string, balance, lockedBalance decimal.Decimal) error {
	if err := ValidateString(walletID, 1, 50); err != nil {
		return err
	}

	if err := validateDecimal(balance); err != nil {
		return err
	}

	return validateDecimal(lockedBalance)
}

func ValidateCreateTradeRequest(BuyerUserEmail, SellerUserEmail, buyOrderId, sellOrderId, marketID string, price, amount, fee decimal.Decimal) error {

	if err := ValidateEmail(BuyerUserEmail); err != nil {
		return err
	}

	if err := ValidateEmail(SellerUserEmail); err != nil {
		return err
	}

	if err := ValidateString(buyOrderId, 1, 50); err != nil {
		return err
	}

	if err := ValidateString(sellOrderId, 1, 50); err != nil {
		return err
	}

	if err := ValidateString(marketID, 1, 50); err != nil {
		return err
	}

	if err := validateDecimal(price); err != nil {
		return err
	}

	if err := validateDecimal(amount); err != nil {
		return err
	}

	if err := validateDecimal(fee); err != nil {
		return err
	}

	return nil
}

func ValidateCreateTransactionRequest(userEmail string, amount decimal.Decimal, transactionType pb.TransactionType) error {
	if err := ValidateEmail(userEmail); err != nil {
		return err
	}

	if err := validateDecimal(amount); err != nil {
		return err
	}

	if transactionType < 0 || transactionType > 2 {
		return fmt.Errorf("transaction_type must be 0 (DEPOSIT), 1 (WITHDRAWAL) or 2 (NONE)")
	}

	return nil
}

func ValidateUpdateTransactionStatusRequest(transactionID string, transactionStatus pb.TransactionStatus) error {
	if err := ValidateString(transactionID, 1, 50); err != nil {
		return err
	}

	if transactionStatus < 0 || transactionStatus > 2 {
		return fmt.Errorf("transaction_status must be 0 (PENDING), 1 (COMPLETED) or 2 (FAILED)")
	}

	return nil
}

func ValidateUpdateOrderStatusAndFilledAmount(orderID string, status pb.Status, filledAmount decimal.Decimal) error {
	if err := ValidateString(orderID, 1, 50); err != nil {
		return err
	}

	if status < 0 || status > 3 {
		return fmt.Errorf("status must be 0 (OPEN), 1 (PARTIALLY_FILLED), 2 (FILLED) or 3 (CANCELLED)")
	}

	if err := validateDecimal(filledAmount); err != nil {
		return err
	}

	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, 32, 128)
}

func validateDecimal(value decimal.Decimal) error {
    coeff := value.Coefficient()          
    precision := len(coeff.String())      
    scale := -value.Exponent()             

    if precision > 20 || scale > 8 {
        return fmt.Errorf("the value exceeds precision of 20 and scale of 8 constraints")
    }
    return nil
}
