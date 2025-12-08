public class TestObject {
    public static void main(String[] args) {
        Counter c = new Counter();
        c.value = 10;
        int x = c.value;  // x = 10
    }
}

class Counter {
    int value;
}