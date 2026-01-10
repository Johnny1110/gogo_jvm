public class TestStringArray {

    public static void main(String[] args) {
        String[] stringArray = {"apple", "banana", "cherry", "date", "elderberry"};

        // Print the original array
        System.out.println("Original array:");
        for (String fruit : stringArray) {
            System.out.println(fruit);
        }

        Class arrayClass = stringArray.getClass();
        System.out.println(arrayClass.getName());
    }
}