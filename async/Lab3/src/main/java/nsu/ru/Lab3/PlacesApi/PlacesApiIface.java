package nsu.ru.Lab3.PlacesApi;

import java.io.IOException;

import nsu.ru.Lab3.controllers.PlacesDTO;

public interface PlacesApiIface {
    PlacesDTO fetchPlacesInRadius(String lat, String lon, String radius) throws IOException, InterruptedException;
    PlaceInfo fetchPlaceDescriptionByXid(String xid)  throws IOException, InterruptedException;
}
