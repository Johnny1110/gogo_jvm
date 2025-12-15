public class TestString {
    public static void main(String[] args) {
        // Test 1: basic print string
        String s1 = "Hello";
        String s2 = "World";
        System.out.println(s1);
        System.out.println(s2);

        // Test 2: 直接輸出字串字面量
        System.out.println("Done");
        System.out.println("<------------------------->");

        if (s1 == s2) {
            System.out.println("s1 == s2");
        } else {
            System.out.println("s1 != s2");
        }

        String s3 = "World";
        if (s2 == s3) {
            System.out.println("s2 == s3");
        } else {
            System.out.println("s2 != s3");
        }
        System.out.println("<------------------------->");

        String s = "你好世界";
        System.out.println(s);
        String mixed = "Hello你好";
        System.out.println(mixed);
        String special = "日本語テスト";
        System.out.println(special);

    }
}