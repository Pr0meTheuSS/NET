package nsu.ru.Lab3.PlacesApi;


import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.concurrent.CompletableFuture;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import nsu.ru.Lab3.Configs.Configs;
import nsu.ru.Lab3.controllers.PlacesDTO;

@Service
public class PlacesApiImpl implements PlacesApiIface {
    private final HttpClient httpClient;
    private String descriptionByXidUrl = "https://api.opentripmap.com/0.1/ru/places/xid/{xid}?apikey={apikey}";
    private String fetchPlacesInRadiusUrl ="http://api.opentripmap.com/0.1/ru/places/radius?radius={radius}&lat={lat}&lon={lon}&format=geojson&apikey={apikey}";
    private String apikey;

    @Autowired
    public PlacesApiImpl(Configs cnfgs) {
        apikey = cnfgs.getPlacesApiKey();
        this.httpClient = HttpClient.newHttpClient();
    }

    @Override
    public CompletableFuture<PlaceInfo> fetchPlaceDescriptionByXid(String xid) throws IOException, InterruptedException {
        String url = descriptionByXidUrl;
        url = url.replace("{xid}", xid);
        url = url.replace("{apikey}", apikey);
        System.out.println(url);
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();
        
        return httpClient.sendAsync(request, HttpResponse.BodyHandlers.ofString())
                    .thenApply(resp -> {
                        ObjectMapper objectMapper = new ObjectMapper();
                        try {
                            return objectMapper.readValue(resp.body(), PlaceInfo.class);
                        } catch (JsonProcessingException e) {
                            e.printStackTrace();
                            return null;
                        }
                    });
    }

    @Override
    public CompletableFuture<PlacesDTO> fetchPlacesInRadius(String lat, String lon, String radius) throws IOException, InterruptedException{
        String url = fetchPlacesInRadiusUrl;
        url = url.replace("{lat}", lat);
        url = url.replace("{lon}", lon);
        url = url.replace("{radius}", radius);
        url = url.replace("{apikey}", apikey);
        System.out.println(url);

        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();
        
        return httpClient.sendAsync(request, HttpResponse.BodyHandlers.ofString())
                    .thenApply(resp -> {
                        ObjectMapper objectMapper = new ObjectMapper();
                        try {
                            return objectMapper.readValue(resp.body(), PlacesDTO.class);
                        } catch (JsonProcessingException e) {
                            e.printStackTrace();
                            return null;
                        }
                    });
    }
}
