public class TestInstanceOfAndCheckCase extends TestParent {

    public TestInstanceOfAndCheckCase(int initialValue) {
        super(initialValue);
    }

    public static void main(String[] args) {
        TestInstanceOfAndCheckCase subClassInstance = new TestInstanceOfAndCheckCase(10);
        System.out.println(subClassInstance instanceof TestParent); // Should print true
        // test case

        Object obj = new TestInstanceOfAndCheckCase(1000);
        TestParent parentInstance = (TestParent) obj;
        parentInstance.printValue(); // Should print 1000
    }
}