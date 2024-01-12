package nsu.ru.Lab3.LocationApi;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.concurrent.CompletableFuture;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.net.URLEncoder;

import nsu.ru.Lab3.Configs.Configs;

@Service
public class LocationApiImpl implements LocationApiIface {
    private final HttpClient httpClient;
    private String fetchLocationsUrl = "https://graphhopper.com/api/1/geocode?q={locationName}&locale=ru&key={apikey}";
    private String apikey;

    @Autowired
    public LocationApiImpl(Configs cnfgs) {
        apikey = cnfgs.getLocationsApiKey();
        this.httpClient = HttpClient.newHttpClient();
    }
    
    public CompletableFuture<LocationResponseDTO> fetchLocations(String locationName) throws IOException, InterruptedException {
        String url = prepareUrlForfetchingLocation(locationName);
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();

        return httpClient.sendAsync(request, HttpResponse.BodyHandlers.ofString())
                .thenApply(x -> {
                    ObjectMapper objectMapper = new ObjectMapper();
                    try {
                        return objectMapper.readValue(x.body(), LocationResponseDTO.class);
                    } catch (JsonMappingException e) {
                        e.printStackTrace();
                        return null;
                    } catch (JsonProcessingException e) {
                        e.printStackTrace();
                        return null;
                    }
                });
    }

    private String prepareUrlForfetchingLocation(String locationName) throws UnsupportedEncodingException {
        String encodedLocationName = URLEncoder.encode(locationName, "UTF-8");
        String url = fetchLocationsUrl;
        url = url.replace("{locationName}", encodedLocationName);
        url = url.replace("{apikey}", apikey);
        return url;
    }
}
