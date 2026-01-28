package java.lang;

/**
 * Thrown to indicate that the {@code clone} method in class
 * {@code Object} has been called to clone an object, but that
 * the object's class does not implement the {@code Cloneable}
 * interface.
 * <p>
 * Applications that override the {@code clone} method can also
 * throw this exception to indicate that an object could not or
 * should not be cloned.
 *
 * @see     java.lang.Cloneable
 * @see     java.lang.Object#clone()
 * @since   1.0
 */

public class CloneNotSupportedException extends Exception {

    private static final long serialVersionUID = 5195511250079656443L;

    /**
     * Constructs a {@code CloneNotSupportedException} with no
     * detail message.
     */
    public CloneNotSupportedException() {
        super();
    }

    /**
     * Constructs a {@code CloneNotSupportedException} with the
     * specified detail message.
     *
     * @param   s   the detail message.
     */
    public CloneNotSupportedException(String s) {
        super(s);
    }
}
