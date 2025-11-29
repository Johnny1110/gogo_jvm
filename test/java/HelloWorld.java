public class HelloWorld {
    // static constant
    private static final String GREETING = "Hello, JVM!";

    // instance
    private int count;

    // constructor
    public HelloWorld() {
        this.count = 0;
    }

    // main
    public static void main(String[] args) {
        System.out.println(GREETING);

        HelloWorld hw = new HelloWorld();
        hw.increment();
        hw.printCount();
    }

    // func 1
    public void increment() {
        count++;
    }

    // func 2
    public void printCount() {
        System.out.println("Count: " + count);
    }

    // static method
    public static int add(int a, int b) {
        return a + b;
    }
}