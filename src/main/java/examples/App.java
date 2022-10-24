package examples;

import org.graalvm.nativeimage.IsolateThread;
import org.graalvm.nativeimage.c.function.CEntryPoint;
import org.graalvm.nativeimage.c.type.CCharPointer;
import org.graalvm.nativeimage.c.type.CIntPointer;
import org.graalvm.nativeimage.c.type.CTypeConversion;
import org.graalvm.nativeimage.c.type.VoidPointer;
import org.graalvm.nativeimage.c.function.CFunctionPointer;

import java.nio.ByteBuffer;
import org.graalvm.nativeimage.c.function.InvokeCFunctionPointer;


interface AllocatorFn extends CFunctionPointer
{
  @InvokeCFunctionPointer
  CCharPointer call(long size);
}
public class App 
{
    @CEntryPoint(name = "pass_bytes")
    public static CCharPointer printStruct(IsolateThread thread, AllocatorFn allocatorFn, VoidPointer message, CIntPointer inSize, CIntPointer outSize) {
        ByteBuffer buf = CTypeConversion.asByteBuffer(message, inSize.read());
        byte[] inArray = new byte[buf.limit()];
        buf.get(inArray);

        System.out.println("In array:");
        for(int i=0; i< inArray.length ; i++) {
          System.out.print(inArray[i] + " ");
       }
        System.out.println();

        byte[] outArray = new byte[3];
        outArray[0] = 1;
        outArray[2] = 2;
        CCharPointer out = allocateAndWrite(allocatorFn, outArray);
        outSize.write(3);
        return out;
    }

    private static CCharPointer allocateAndWrite(AllocatorFn allocatorFn, byte[] b)
    {
      CCharPointer a = allocatorFn.call(b.length + 1);
      for (int i = 0; i < b.length; i++) {
        a.write(i, b[i]);
      }
      return a;
    }
}
