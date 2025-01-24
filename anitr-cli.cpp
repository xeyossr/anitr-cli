#include "modules/anitr_fetch.h"
#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <cstdlib>
#include <fstream>
#include <sstream>
#include <algorithm>

FetchData fetchdata;

// Ana menü seçenekleri
std::vector<std::string> main_menu_options = {
	"İzle", "Sonraki Bölüm", "Önceki Bölüm", "Bölüm Seç", "Anime Ara", "Çık"
};

std::string movie_url;
bool is_movie;

// Yardım menüsü
void printHelp() 
{
    std::cout << "anitr-cli kullanımı:\n"
              << "  --help, -h: Bu yardım menüsünü gösterir\n"
              << "  --gen-config: rofi-flags.conf dosyasını oluşturur\n"
              << "\n";
}

// rofi.flags.conf dosyasını oluşturacak funksiyon
void generateConfigFile() 
{
    std::string configDir = std::string(getenv("HOME")) + "/.config/anitr-cli";
    std::string configFile = configDir + "/rofi-flags.conf";

    // Klasörü oluştur (varsa atla)
    std::filesystem::create_directories(configDir);

    // Konfigürasyon dosyası varsa, uyarı göster ve çık
    if (std::filesystem::exists(configFile)) 
    {
        std::cout << "Konfigürasyon dosyası zaten var: " << configFile << "\n";
        return;
    }

    // Varsayılan parametrelerle dosya oluştur
    std::ofstream config(configFile);
    if (config.is_open()) 
    {
        config << "";
        config.close();
        std::cout << "Konfigürasyon dosyası oluşturuldu: " << configFile << "\n";
    } 
    
    else 
    {
        std::cerr << "Konfigürasyon dosyası oluşturulamıyor: " << configFile << "\n";
    }
}


// Kullanıcıdan rofi ile giriş alacak fonksiyon
std::string getInputFromRofi(const std::string& prompt, const std::vector<std::string>& options) 
{
    // Rofi parametrelerini okumak için konfigürasyon dosyasını kontrol et
    std::string rofi_flags = "";
    std::ifstream config_file(std::string(getenv("HOME")) + "/.config/anitr-cli/rofi-flags.conf");

    if (config_file.is_open()) 
    {
        std::string line;
        while (std::getline(config_file, line)) 
        {
            // Eğer satır boş değilse ve ya başında # yoksa
            if (!line.empty() && line[0] != '#') 
            {
                rofi_flags += line + " ";  // Parametreleri birleştir
            }
        }
        
        config_file.close();
    }

    // Rofi komutunu oluştur
    std::string rofi_cmd = "{ echo '" + prompt + "\n" + 
                            "\n" + 
                            "[\n" + 
                            "  '" + prompt + "'\n" +
                            "'<back>'\n" +
                            "'<exit>'" + 
                            "']\n\n" ;

    rofi_cmd += "echo -e \"" ;

    for (const auto& option : options) 
    {
        rofi_cmd += option + "\n";
    }

    rofi_cmd += "\" | rofi -dmenu -p '" + prompt + "' " + rofi_flags + "; } 2>/dev/null";  // Çıktıyı gizle

    std::string selected;
    
    FILE* fp = popen(rofi_cmd.c_str(), "r");
    if (fp != NULL) 
    {
        char buffer[1024];
        
        if (fgets(buffer, sizeof(buffer), fp) != NULL) 
        {
            selected = buffer;
        }
        
        fclose(fp);
    }

    // Satır sonu karakterini kaldır
    if (!selected.empty() && selected.back() == '\n') 
    {
        selected.pop_back();
    }

    return selected;
}


std::vector<std::map<std::string, std::string>> queryLoop() 
{
        std::string query;
        std::vector<std::map<std::string, std::string>> results;

        // Anime arama prompt'unu döngüye al
        while (true) 
        {
            query = getInputFromRofi("Anime Ara", {"Çık"});

            if (query == "<exit>" || query == "exit" || query == "Çık") 
            {
                exit(0);
            }

            if (query.find_first_of("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") == std::string::npos) 
            {
                continue;
            }

            results = fetchdata.fetch_anime_search_data(query);

            if (results.empty()) 
            {
                std::cout << "Sonuç bulunamadı. Tekrar deneyin!" << "\n";
                continue;
            }

            break;
        }

        return results;
}


int main(int argc, char* argv[]) {

    // Komut satırı parametrelerini kontrol et
    for (int i = 1; i < argc; ++i) 
    {
        std::string arg = argv[i];

        if (arg == "--help" || arg == "-h") 
        {
            printHelp();
            return 0;
        } 
        
        else if (arg == "--gen-config") 
        {
            generateConfigFile();
            return 0;
        }
    }

    std::vector<std::map<std::string, std::string>> anime_episodes;
    int selected_episode_index = 0;

    // Anime arama prompt'u
    std::string query;

    // Arama sonuçlarını al
    std::vector<std::map<std::string, std::string>> results = queryLoop();
    

    std::vector<std::string> anime_names;
    for (const auto& item : results) 
    {
        anime_names.push_back(item.at("name"));
    }

    anime_names.push_back("Çık");

    std::string selected_anime_name;

    while (true) 
    {
        // Anime seçimi
        selected_anime_name = getInputFromRofi("Anime Seç", anime_names);

        // Eğer arama kısmına <exit> ya da exit yazılırsa çık
        if (selected_anime_name == "<exit>" || selected_anime_name == "exit") exit(0);

        // Şartlar
        bool is_not_in_anime_names = std::find(anime_names.begin(), anime_names.end(), selected_anime_name) == anime_names.end();
        bool is_not_exit_command = selected_anime_name != "Çık";
        
        // Eğer selected_anime_name, anime_names içerisindeki bir öğe değilse ve Çık değilse döngüye devam et
        if (is_not_in_anime_names && is_not_exit_command) 
        {
            continue;
        }

        // Döngüden çık
        break;

    }

    // Seçilen animeyi bul
    std::map<std::string, std::string> selected_anime;
    for (const auto& item : results) 
    {
        if (item.at("name") == selected_anime_name) 
        {
            selected_anime = item;
            break;
        }
    }

    // Seçilen animenin bölümlerini al
    std::string selected_id = selected_anime.at("id");
    anime_episodes = fetchdata.fetch_anime_episodes(selected_id);

    // Eğer animede herhangi bir bölüm bulunamadıysa
    is_movie = anime_episodes.empty();
    
    // Eğer film seçildiyse
    if (is_movie) 
    {
        movie_url = fetchdata.fetch_anime_watch_api_url_movie(selected_id);
    }

    while (true) 
    {
		// Ana menüye bölümü de ekle
        if (!is_movie)
        {
            main_menu_options = {
			    "İzle", "Sonraki Bölüm", "Önceki Bölüm", "Bölüm Seç", "Anime Ara", "Çık", anime_episodes[selected_episode_index].at("name")
		    };
        }

        else if (is_movie) {
            main_menu_options = {
                "İzle", "Anime Ara", "Çık"
            };
        }

        // Ana menüyü göster
        std::string main_menu_choice = getInputFromRofi("Ana Menü", main_menu_options);

        // Eğer çık seçeneği seçildiyse    
        if (main_menu_choice == "Çık") 
        {
            exit(0);
            break;
        } 
        
        // Eğer anime ara seçeneği seçildiyse
        else if (main_menu_choice == "Anime Ara") 
        {
<<<<<<< HEAD

			// Selected Episode Index'i sıfırla
        	selected_episode_index = 0;
        	
=======
	    // episode_index'i sıfırla
	    selected_episode_index = 0;
		
>>>>>>> 2ce9a9d8666f18a0f4da0e77d40dfdca9b0670f8
            // Anime Ara işlemini tekrar başlatıyoruz
            std::vector<std::map<std::string, std::string>> results = queryLoop();

            // Yeni sonuçları anime_names listesine ekleyin
            anime_names.clear();  // Önceki verileri temizle

            for (const auto& item : results) 
            {
                anime_names.push_back(item.at("name"));
            }
            
            anime_names.push_back("Çık");

            // Anime seçme menüsünü tekrar göster
            selected_anime_name = getInputFromRofi("Anime Seç", anime_names);

            // Eğer arama kısmına <exit> ya da exit yazılırsa çık
            if (selected_anime_name == "<exit>" || selected_anime_name == "exit") 
            {
                exit(0);
            }

            // Şartlar
            bool is_not_in_anime_names = std::find(anime_names.begin(), anime_names.end(), selected_anime_name) == anime_names.end();
            bool is_not_exit_command = selected_anime_name != "Çık";

            // Eğer selected_anime_name, anime_names içerisindeki bir öğe değilse ve Çık değilse döngüye devam et
            if (is_not_in_anime_names && is_not_exit_command) 
            {
                continue;
            }

            // Seçilen animeyi bul
            std::map<std::string, std::string> selected_anime;
            for (const auto& item : results) 
            {
                if (item.at("name") == selected_anime_name) 
                {
                    selected_anime = item;
                    break;
                }
            }

            // Seçilen animenin bölümlerini al
            std::string selected_id = selected_anime.at("id");
            anime_episodes = fetchdata.fetch_anime_episodes(selected_id);

            // Eğer animede herhangi bir bölüm bulunamadıysa
            is_movie = anime_episodes.empty();

            // Eğer film seçildiyse
            if (is_movie) 
            {
                movie_url = fetchdata.fetch_anime_watch_api_url_movie(selected_id);

                /*
                if (!movie_url.empty()) 
                {
                    std::string mpv_cmd = "mpv --fullscreen " + movie_url + " > /dev/null 2>&1";
                    std::cout << "İzleniyor: " << selected_anime_name << "\n";
                    system(mpv_cmd.c_str());
                } 

                else 
                    std::cout << "Filmin URL'si bulunamadı." << "\n";
                */

                return 0;
            }
        }
        
        // Eğer izle seçeneği seçildiyse
        else if (main_menu_choice == "İzle") 
        {

            // Eğer film seçilmediyse

			if (!is_movie) 
            {
            	// Seçilen bölümün URL'sini al
            	std::string episode_url = anime_episodes[selected_episode_index].at("url");
            
            	// Bölüm URL'si ile izleme URL'sini al
            	std::vector<std::map<std::string, std::string>> watch_url = fetchdata.fetch_anime_watch_api_url(episode_url);

                if (!watch_url.empty()) 
                {
                    // URL'yi al
                    std::string video_url = watch_url.back().at("url");

                    // MPV ile izleme başlat
                    std::cout << "İzleniyor: " << selected_anime_name << " " << anime_episodes[selected_episode_index].at("name") << "\n";
                    std::string mpv_cmd = "mpv --fullscreen " + video_url + " > /dev/null 2>&1";
                    system(mpv_cmd.c_str());
                } 
            
                else 
                {
                    std::cerr << "İzleme URL'si alınamadı" << "\n";
                }
          
            } 
          
            else if (is_movie)
            {
                if (!movie_url.empty()) 
                {
          	    	// MPV ile izleme başlat
          	    	std::cout << "İzleniyor: " << selected_anime_name << "\n";
          	    	std::string mpv_cmd = "mpv --fullscreen " + movie_url + " > /dev/null 2>&1";
          	    	system(mpv_cmd.c_str());
          	    } 
            
                else 
                {
          	    	std::cerr << "İzleme URL'si alınamadı" << "\n";
          	    }
            }    
        
        } 
        
        // Eğer sonraki bölüm seçeneği seçildiyse
        else if (main_menu_choice == "Sonraki Bölüm") 
        {
            if (selected_episode_index < anime_episodes.size() - 1) 
            {
                selected_episode_index++;
                
                std::string episode_url = anime_episodes[selected_episode_index].at("url");
                std::vector<std::map<std::string, std::string>> watch_url = fetchdata.fetch_anime_watch_api_url(episode_url);

                // URL'yi al
                std::string video_url = watch_url.back().at("url");

                // MPV ile izleme başlat
                std::cout << "İzleniyor: " << selected_anime_name << " " << anime_episodes[selected_episode_index].at("name") << "\n";
                std::string mpv_cmd = "mpv --fullscreen " + video_url + " > /dev/null 2>&1";
                system(mpv_cmd.c_str());

            } 
            
            else 
            {
                std::cout << "Zaten en son bölümdesiniz" << "\n";
            }
        
        } 
        
        // Eğer önceki bölüm seçeneği seçildiyse
        else if (main_menu_choice == "Önceki Bölüm") 
        {
            if (selected_episode_index > 0) 
            {
                selected_episode_index--;

                std::string episode_url = anime_episodes[selected_episode_index].at("url");
                std::vector<std::map<std::string, std::string>> watch_url = fetchdata.fetch_anime_watch_api_url(episode_url);
                
                // URL'yi al
                std::string video_url = watch_url.back().at("url");

                // MPV ile izleme başlat
                std::cout << "İzleniyor: " << selected_anime_name << " " << anime_episodes[selected_episode_index].at("name") << "\n";
                std::string mpv_cmd = "mpv --fullscreen " + video_url + " > /dev/null 2>&1";
                system(mpv_cmd.c_str());

            } 
            
            else 
            {
                std::cout << "Zaten ilk bölümdesiniz" << "\n";
            }

        } 
        
        // Eğer bölüm seç seçeneği seçildiyse
        else if (main_menu_choice == "Bölüm Seç") 
        {
            // Bölüm listesini göster ve kullanıcıyı seçim yapmaya yönlendir
            std::vector<std::string> episode_titles;
            
            for (const auto& episode : anime_episodes) 
            {
                episode_titles.push_back(episode.at("name"));
            }

            std::string selected_episode_title = getInputFromRofi("Bölüm Seç", episode_titles);
            //if (selected_episode_title == "<exit>" || selected_episode_title.empty()) break;

            // Seçilen bölüm verisini bul
            for (int i = 0; i < anime_episodes.size(); i++) 
            {
                if (anime_episodes[i].at("name") == selected_episode_title) 
                {
                    selected_episode_index = i;
                    break;
                }
            }
        }
    }

    return 0;
}
