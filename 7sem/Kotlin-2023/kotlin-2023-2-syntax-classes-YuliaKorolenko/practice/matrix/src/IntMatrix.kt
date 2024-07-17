class IntMatrix(val rows: Int, val columns: Int) {
    private var matrix: IntArray

    init {
        if (rows <= 0 || columns <= 0) {
            throw IllegalArgumentException("Can't create matrix with negative rows or columns")
        }
        this.matrix = IntArray(rows * columns)
    }

    operator fun get(
        i: Int,
        j: Int,
    ): Any {
        checkValue(i, j)
        return matrix[i * columns + columns]
    }

    operator fun set(
        i: Int,
        j: Int,
        value: Int,
    ) {
        checkValue(i, j)
        matrix[i * columns + columns] = value
    }

    @Throws(IllegalArgumentException::class)
    private fun checkValue(
        i: Int,
        j: Int,
    ) {
        if (i !in 0..rows || j !in 0..columns) {
            throw IllegalArgumentException("Don't exist this field")
        }
    }
}
