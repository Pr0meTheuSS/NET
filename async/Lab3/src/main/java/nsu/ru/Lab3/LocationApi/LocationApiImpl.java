package nsu.ru.Lab3.LocationApi;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.URLEncoder;

import com.fasterxml.jackson.databind.ObjectMapper;


public class LocationApiImpl implements LocationApiIface {
        private final HttpClient httpClient;
    
    public LocationApiImpl() {
        this.httpClient = HttpClient.newHttpClient();
    }

    @Override
    public LocationResponseDTO fetchLocations(String locationName) throws IOException, InterruptedException {
        ObjectMapper objectMapper = new ObjectMapper();
        String ret;
        try {
            ret = fetchLocation(locationName);
            return objectMapper.readValue(ret, LocationResponseDTO.class);
        } catch (IOException | InterruptedException e) {
            throw e;
        }
    }
    
    private String fetchLocation(String locationName) throws IOException, InterruptedException {
        String url = prepareUrlForfetchingLocation(locationName);
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .build();

        HttpResponse<String> resp;
        resp = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return resp.body();
    }

    private String prepareUrlForfetchingLocation(String locationName) throws UnsupportedEncodingException {
        String encodedValue = "";
        encodedValue = URLEncoder.encode(locationName, "UTF-8");
        System.out.println(encodedValue);

        return "https://graphhopper.com/api/1/geocode?q=" + encodedValue + "&locale=ru&key=b21dcab5-cc27-472b-8a3e-d1eb62c38a04";
    }
}
