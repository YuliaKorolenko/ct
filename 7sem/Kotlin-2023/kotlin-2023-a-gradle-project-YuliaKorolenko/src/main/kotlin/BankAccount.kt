class BankAccount(amount: Int) {
    var balance: Int = 0
        private set(newValue) {
            logTransaction(balance, newValue)
            field = newValue
        }

    init {
        balance =
            if (!checkNegative(amount)) {
                amount
            } else {
                throw IllegalArgumentException("Can't create account with negative balance")
            }
    }

    private fun checkNegative(value: Int): Boolean = value <= 0

    fun withdraw(value: Int) {
        if (balance < value || checkNegative(value)) {
            throw IllegalArgumentException("Illegal withdrawal amount")
        }
        balance -= value
    }

    fun deposit(value: Int) {
        if (checkNegative(value)) {
            throw IllegalArgumentException("Illegal replenishment amount")
        }
        balance += value
    }
}
