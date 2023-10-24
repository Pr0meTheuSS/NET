package nsu.ru.Lab3.PlacesApi;


import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

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
    public PlaceInfo fetchPlaceDescriptionByXid(String xid) throws IOException, InterruptedException {
        String url = descriptionByXidUrl;
        url = url.replace("{xid}", xid);
        url = url.replace("{apikey}", apikey);
        System.out.println(url);
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
        String url = fetchPlacesInRadiusUrl;
        url = url.replace("{lat}", lat);
        url = url.replace("{lon}", lon);
        url = url.replace("{radius}", radius);
        url = url.replace("{apikey}", apikey);
        System.out.println(url);

        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();
        
        HttpResponse<String> resp;
        resp = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

        ObjectMapper objectMapper = new ObjectMapper();
        return objectMapper.readValue(resp.body(), PlacesDTO.class);
    }
}
