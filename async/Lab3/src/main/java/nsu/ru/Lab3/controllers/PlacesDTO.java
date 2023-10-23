package nsu.ru.Lab3.controllers;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public class PlacesDTO {
    private String type;
    private List<Feature> features;

    @JsonProperty("type")
    public String getType() {
        return type;
    }

    @JsonProperty("type")
    public void setType(String type) {
        this.type = type;
    }

    @JsonProperty("features")
    public List<Feature> getFeatures() {
        return features;
    }

    @JsonProperty("features")
    public void setFeatures(List<Feature> features) {
        this.features = features;
    }
}

@JsonIgnoreProperties(ignoreUnknown = true)
class Feature {
    private String type;
    private String id;
    private Geometry geometry;
    private Properties properties;

    @JsonProperty("type")
    public String getType() {
        return type;
    }

    @JsonProperty("type")
    public void setType(String type) {
        this.type = type;
    }

    @JsonProperty("id")
    public String getId() {
        return id;
    }

    @JsonProperty("id")
    public void setId(String id) {
        this.id = id;
    }

    @JsonProperty("geometry")
    public Geometry getGeometry() {
        return geometry;
    }

    @JsonProperty("geometry")
    public void setGeometry(Geometry geometry) {
        this.geometry = geometry;
    }

    @JsonProperty("properties")
    public Properties getProperties() {
        return properties;
    }

    @JsonProperty("properties")
    public void setProperties(Properties properties) {
        this.properties = properties;
    }
}

@JsonIgnoreProperties(ignoreUnknown = true)
class Geometry {
    private String type;
    private List<Double> coordinates;

    @JsonProperty("type")
    public String getType() {
        return type;
    }

    @JsonProperty("type")
    public void setType(String type) {
        this.type = type;
    }

    @JsonProperty("coordinates")
    public List<Double> getCoordinates() {
        return coordinates;
    }

    @JsonProperty("coordinates")
    public void setCoordinates(List<Double> coordinates) {
        this.coordinates = coordinates;
    }
}

@JsonIgnoreProperties(ignoreUnknown = true)
class Properties {
    private String xid;
    private String name;
    private double dist;
    private String kinds;

    @JsonProperty("xid")
    public String getXid() {
        return xid;
    }

    @JsonProperty("xid")
    public void setXid(String xid) {
        this.xid = xid;
    }

    @JsonProperty("name")
    public String getName() {
        return name;
    }

    @JsonProperty("name")
    public void setName(String name) {
        this.name = name;
    }

    @JsonProperty("dist")
    public double getDist() {
        return dist;
    }

    @JsonProperty("dist")
    public void setDist(double dist) {
        this.dist = dist;
    }

    @JsonProperty("kinds")
    public String getKinds() {
        return kinds;
    }

    @JsonProperty("kinds")
    public void setKinds(String kinds) {
        this.kinds = kinds;
    }
}
