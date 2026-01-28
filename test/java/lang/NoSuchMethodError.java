package java.lang;

/**
 * Thrown if an application tries to call a specified method of a
 * class (either static or instance), and that class no longer has a
 * definition of that method.
 * <p>
 * Normally, this error is caught by the compiler; this error can
 * only occur at run time if the definition of a class has
 * incompatibly changed.
 *
 * @since   1.0
 */
public class NoSuchMethodError {
    private static final long serialVersionUID = -3765521442372831335L;

    /**
     * Constructs a {@code NoSuchMethodError} with no detail message.
     */
    public NoSuchMethodError() {
    }

    /**
     * Constructs a {@code NoSuchMethodError} with the
     * specified detail message.
     *
     * @param   s   the detail message.
     */
    public NoSuchMethodError(String s) {
    }
}
