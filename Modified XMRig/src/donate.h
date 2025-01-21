/* XMRig
 * Copyright (c) 2018-2022 SChernykh   <https://github.com/SChernykh>
 * Copyright (c) 2016-2022 XMRig       <https://github.com/xmrig>, <support@xmrig.com>
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef XMRIG_DONATE_H
#define XMRIG_DONATE_H

 /* Fallen Miner
  * 
  * You should of downloaded this code from https://github.com/laprosa/Fallen-Miner, if you haven't please examine ALL files to ensure they have not been tampered with.
  * The current donation level is set to 2%, which is 1% higher than the default configuration towards the xmrig project.
  * Example of how it works for the setting of 2%:
  * Your miner will mine into your usual pool for a random time (in a range from 49.5 to 148.5 minutes),
  * then switch to the developer's pool for 2 minutes, then switch again to your pool for 98 minutes
  * and then switch again to developer's pool for 1 minute; these rounds will continue until the miner stops.
  * If you don't want to donate for 2 minutes via hashpower, set both values to 0.
  * If you don't donate I completely understand however as I have worked for a while on this project, please consider a one off donation to the addreses in the project readme.
  * Thank you for using my software :)
  */


constexpr const int kDefaultDonateLevel = 2;
constexpr const int kMinimumDonateLevel = 2;


#endif // XMRIG_DONATE_H
