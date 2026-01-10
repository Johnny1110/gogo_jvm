package java.lang;

public class Class {

    // Cache the name to reduce the number of calls into the VM.
    // This field would be set by VM itself during initClassName call.
    private transient String name;

    public String getName() {
            String name = this.name;
            return name != null ? name : initClassName();
    }

    private native String initClassName();
   
}