package nsu.ru.Lab3.controllers;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;

public class LocationResponseDTO {
    @JsonProperty("hits")
    private Location[] hits;

    @JsonProperty("locale")
    private String locale;

    public List <Location> getHits() {
        return List.of(hits);
    }

    public void setHits(Location[] hits) {
        this.hits = hits;
    }

    public String getLocale() {
        return locale;
    }

    public void setLocale(String locale) {
        this.locale = locale;
    }
}

class Location {
    @JsonIgnore
    private Point point;

    public Point getPoint() {
        return point;
    }
    public void setPoint(Point point) {
        this.point = point;
    }
    
    private double[] extent;
    private String name;
    private String country;
    private String countrycode;
    private String city;
    private String state;
    private String postcode;

    private String street;
    private String housenumber;

    public String getStreet() {
        return street;
    }
    public void setStreet(String street) {
        this.street = street;
    }
    public String getHousenumber() {
        return housenumber;
    }
    public void setHousenumber(String housenumber) {
        this.housenumber = housenumber;
    }

    @JsonIgnore
    private long osm_id;
    @JsonIgnore
    private String osm_type;
    @JsonIgnore
    private String osm_key;
    @JsonIgnore
    private String osm_value;
    public long getOsm_id() {
        return osm_id;
    }
    public void setOsm_id(long osm_id) {
        this.osm_id = osm_id;
    }
    public String getOsm_type() {
        return osm_type;
    }
    public void setOsm_type(String osm_type) {
        this.osm_type = osm_type;
    }
    public String getOsm_key() {
        return osm_key;
    }
    public void setOsm_key(String osm_key) {
        this.osm_key = osm_key;
    }
    public String getOsm_value() {
        return osm_value;
    }
    public void setOsm_value(String osm_value) {
        this.osm_value = osm_value;
    }
    public double[] getExtent() {
        return extent;
    }
    public void setExtent(double[] extent) {
        this.extent = extent;
    }
    public String getName() {
        return name;
    }
    public void setName(String name) {
        this.name = name;
    }
    public String getCountry() {
        return country;
    }
    public void setCountry(String country) {
        this.country = country;
    }
    public String getCountrycode() {
        return countrycode;
    }
    public void setCountrycode(String countrycode) {
        this.countrycode = countrycode;
    }
    public String getCity() {
        return city;
    }
    public void setCity(String city) {
        this.city = city;
    }
    public String getState() {
        return state;
    }
    public void setState(String state) {
        this.state = state;
    }
    public String getPostcode() {
        return postcode;
    }
    public void setPostcode(String postcode) {
        this.postcode = postcode;
    }

}

class Point {
    private String lat;
    private String lng;

    public String getLat() {
        return lat;
    }
    public void setLat(String lat) {
        this.lat = lat;
    }
    public String getLng() {
        return lng;
    }
    public void setLng(String lng) {
        this.lng = lng;
    }
}