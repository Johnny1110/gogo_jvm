public class TestSubClass extends TestParent {

    public TestSubClass(int initialValue) {
        super(initialValue);
    }

    public void display() {
        printValue();
    }

    public void setValue(int newValue) {
        this.value = newValue;
    }

    public static void main(String[] args) {
        TestSubClass subClassInstance = new TestSubClass(10);
        subClassInstance.display(); // Should print 10
        subClassInstance.setValue(20);
        subClassInstance.display(); // Should print 20
    }
}
