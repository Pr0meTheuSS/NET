package nsu.ru.Lab3.LocationApi;
import java.io.IOException;

public interface LocationApiIface {
    LocationResponseDTO fetchLocations(String locationName) throws IOException, InterruptedException;
}
