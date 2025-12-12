public class TestPrintlnCalc {
    public static void main(String[] args) {
        int a = 10;
        int b = 20;
        System.out.println(a + b);  // 30

        System.out.println(fib(10));  // 55
    }

    public static int fib(int n) {
        if (n <= 1) return n;
        return fib(n - 1) + fib(n - 2);
    }
}