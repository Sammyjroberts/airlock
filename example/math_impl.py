class MathImplementation:
    def Add(self, a: int, b: int) -> int:
        return a + b

    def Subtract(self, a: int, b: int) -> int:
        return a - b

    def Multiply(self, a: int, b: int) -> int:
        return a * b

    def Divide(self, a: int, b: int) -> int:
        if b == 0:
            raise ValueError("division by zero")
        return a // b
