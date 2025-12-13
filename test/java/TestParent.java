public class TestParent {
    protected int value;

    public TestParent(int initialValue) {
        this.value = initialValue;
    }

    protected void printValue() {
        System.out.println(value);
    }
}
