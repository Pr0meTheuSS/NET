package nsu.ru.Lab3.PlacesApi;


import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import com.fasterxml.jackson.databind.ObjectMapper;

import nsu.ru.Lab3.controllers.PlacesDTO;

public class PlacesApiImpl implements PlacesApiIface {
    private final HttpClient httpClient;

    public PlacesApiImpl() {
        this.httpClient = HttpClient.newHttpClient();
    }

    @Override
    public PlaceInfo fetchPlaceDescriptionByXid(String xid) throws IOException, InterruptedException {
        String url = "https://api.opentripmap.com/0.1/ru/places/xid/" + xid + "?apikey=5ae2e3f221c38a28845f05b6571f62f3300e847fa60718705f4aad5d";
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();
        
        HttpResponse<String> resp;
        resp = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        System.out.println(resp.body());

        ObjectMapper objectMapper = new ObjectMapper();
        return objectMapper.readValue(resp.body(), PlaceInfo.class);
    }

    @Override
    public PlacesDTO fetchPlacesInRadius(String lat, String lon, String radius) throws IOException, InterruptedException{
        String url = "http://api.opentripmap.com/0.1/ru/places/radius?radius=" + radius + "&lat=" + lat + 
        "&lon=" + lon + "&format=geojson&apikey=5ae2e3f221c38a28845f05b6571f62f3300e847fa60718705f4aad5d";
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();
        
        HttpResponse<String> resp;
        resp = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

        ObjectMapper objectMapper = new ObjectMapper();
        return objectMapper.readValue(resp.body(), PlacesDTO.class);
    }
}
