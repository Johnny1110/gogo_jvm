public class TestStaticField {
    static int value = 100;

    public static void main(String[] args) {
        int a = value;      // getstatic
        value = 200;        // putstatic
        int b = value;      // getstatic
        // b = 200
    }
}