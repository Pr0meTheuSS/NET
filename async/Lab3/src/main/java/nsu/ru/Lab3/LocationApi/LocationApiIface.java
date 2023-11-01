package nsu.ru.Lab3.LocationApi;
import java.io.IOException;
import java.util.concurrent.CompletableFuture;

public interface LocationApiIface {
    CompletableFuture<LocationResponseDTO> fetchLocations(String locationName) throws IOException, InterruptedException;
}
