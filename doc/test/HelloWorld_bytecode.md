# HelloWorld.java

<br>

## Source Code:

```java
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
```

<br>
<br>

## After Compile

```
###### Class defpackage.HelloWorld (HelloWorld)
// class version 52.0 (52)
// access flags 0x21
public class HelloWorld {

  // compiled from: HelloWorld.java

  // access flags 0x1A
  private final static Ljava/lang/String; GREETING = "Hello, JVM!"

  // access flags 0x2
  private I count

  // access flags 0x1
  public <init>()V
   L0
    LINENUMBER 9 L0
    ALOAD 0
    INVOKESPECIAL java/lang/Object.<init> ()V
   L1
    LINENUMBER 10 L1
    ALOAD 0
    ICONST_0
    PUTFIELD HelloWorld.count : I
   L2
    LINENUMBER 11 L2
    RETURN
    MAXSTACK = 2
    MAXLOCALS = 1

  // access flags 0x9
  public static main([Ljava/lang/String;)V
   L0
    LINENUMBER 15 L0
    GETSTATIC java/lang/System.out : Ljava/io/PrintStream;
    LDC "Hello, JVM!"
    INVOKEVIRTUAL java/io/PrintStream.println (Ljava/lang/String;)V
   L1
    LINENUMBER 17 L1
    NEW HelloWorld
    DUP
    INVOKESPECIAL HelloWorld.<init> ()V
    ASTORE 1
   L2
    LINENUMBER 18 L2
    ALOAD 1
    INVOKEVIRTUAL HelloWorld.increment ()V
   L3
    LINENUMBER 19 L3
    ALOAD 1
    INVOKEVIRTUAL HelloWorld.printCount ()V
   L4
    LINENUMBER 20 L4
    RETURN
    MAXSTACK = 2
    MAXLOCALS = 2

  // access flags 0x1
  public increment()V
   L0
    LINENUMBER 24 L0
    ALOAD 0
    DUP
    GETFIELD HelloWorld.count : I
    ICONST_1
    IADD
    PUTFIELD HelloWorld.count : I
   L1
    LINENUMBER 25 L1
    RETURN
    MAXSTACK = 3
    MAXLOCALS = 1

  // access flags 0x1
  public printCount()V
   L0
    LINENUMBER 29 L0
    GETSTATIC java/lang/System.out : Ljava/io/PrintStream;
    NEW java/lang/StringBuilder
    DUP
    INVOKESPECIAL java/lang/StringBuilder.<init> ()V
    LDC "Count: "
    INVOKEVIRTUAL java/lang/StringBuilder.append (Ljava/lang/String;)Ljava/lang/StringBuilder;
    ALOAD 0
    GETFIELD HelloWorld.count : I
    INVOKEVIRTUAL java/lang/StringBuilder.append (I)Ljava/lang/StringBuilder;
    INVOKEVIRTUAL java/lang/StringBuilder.toString ()Ljava/lang/String;
    INVOKEVIRTUAL java/io/PrintStream.println (Ljava/lang/String;)V
   L1
    LINENUMBER 30 L1
    RETURN
    MAXSTACK = 3
    MAXLOCALS = 1

  // access flags 0x9
  public static add(II)I
   L0
    LINENUMBER 34 L0
    ILOAD 0
    ILOAD 1
    IADD
    IRETURN
    MAXSTACK = 2
    MAXLOCALS = 2
}
```