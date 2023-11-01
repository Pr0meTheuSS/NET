package nsu.ru.Lab3.PlacesApi;

import java.io.IOException;
import java.util.concurrent.CompletableFuture;

import nsu.ru.Lab3.controllers.PlacesDTO;


public interface PlacesApiIface {
    CompletableFuture<PlacesDTO> fetchPlacesInRadius(String lat, String lon, String radius) throws IOException, InterruptedException;
    CompletableFuture<PlaceInfo> fetchPlaceDescriptionByXid(String xid)  throws IOException, InterruptedException;
}
